package master

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	. "mervin.me/SalonRabbit/util/log"
	"mervin.me/util/ip"
)
import (
	"sync"
	"sync/atomic"
	// "time"
	"encoding/json"
	"errors"
	"fmt"
)

type Superstep struct {
	Curstep    uint32
	SignalChan *com.SignalChannel
	slaveState map[uint32]int //map[ipLong]active nodes of the ip address. -1: all the node

	runningStep chan int

	counter  int32
	slaveNum int32
	mutex    *sync.Mutex
}

func NewSuperstep(sc *com.SignalChannel) *Superstep {
	Log.Print("New")
	s := &Superstep{
		Curstep:     0,
		slaveState:  make(map[uint32]int, 100),
		runningStep: make(chan int, 1),
		counter:     0,
		slaveNum:    int32(len(SlaveHosts)),
		mutex:       new(sync.Mutex),
	}
	for _, h := range SlaveHosts {
		s.slaveState[ip.Ip2Int(h)] = -1
	}
	return s
}

//
// traverse all the slave with active nodes
// traverse done
// wait  slave running is done
// request the sync: active node
// response the sync request
// next step
func (s *Superstep) Run() {
	Log.Print("Superstep:Run")
	// go func() { s.runningStep <- 1 }()
	s.runningStep <- 1
	for {
		v, ok := <-s.runningStep
		if !ok {
			break
		}

		switch v {
		case 1:
			//traverse slave
			s.SendNextStep()
		case 2:
			//
			s.SendSyncRequest()
		case 3:
			//is terminate
			if atomic.LoadInt32(&s.counter) > 0 {
				atomic.StoreInt32(&s.counter, 0)
				s.runningStep <- 1
			} else {
				//terminated
				close(s.runningStep)
				break
			}
		}
	}
	println("End:Run")
}

func (s *Superstep) SendNextStep() {
	Log.Print(">>SendNextStep ", s.Curstep+1)
	s.Curstep++
	for dstAddr, activeNodeNum := range s.slaveState {
		if activeNodeNum != 0 {
			ns := message.Nextstep{
				Host: LocalIP,
				Step: s.Curstep,
			}
			var target engine.ByteIndex
			str := fmt.Sprintf("%08x", dstAddr)
			copy(target[0:8], str)
			msg := &message.Message{
				Object: message.LOCAL,
				Action: message.NEXTSTEP,
				Target: target,
				Value:  ns.Marshal(),
			}
			var reply message.Reply
			for {
				err := s.SignalChan.Send(msg, &reply)

				if reply.State == 200 && err == nil {
					break
				}
			}
		}
	}

}

func (s *Superstep) SendSyncRequest() {
	Log.Print(">>SendSyncRequest ", s.Curstep)
	for dstAddr, _ := range s.slaveState {
		var target engine.ByteIndex
		str := fmt.Sprintf("%08x", dstAddr)
		copy(target[0:8], str)
		v, _ := json.Marshal(s.Curstep)
		msg := &message.Message{
			Object: message.LOCAL,
			Action: message.SYNC,
			Target: target,
			Value:  v,
		}
		var reply message.Reply
		for {
			err := s.SignalChan.Send(msg, &reply)
			if reply.State == 200 && err == nil {
				var num int
				json.Unmarshal(reply.Value, &num)
				if num > 0 {
					Log.Print("<<SendSyncRequest Done \t", dstAddr, "\t", num)
					s.slaveState[dstAddr] = num
					atomic.AddInt32(&s.counter, 1)
				}
				break
			}
		}
	}
	s.runningStep <- 3
}
func (s *Superstep) Put(msg *message.Message, reply *message.Reply) (err error) {
	Log.Print("Put")
	if msg.Action == message.SLAVEDONE {
		atomic.AddInt32(&s.counter, 1)
		Log.Print("<<Slave Done: ", s.counter)
		if s.counter == s.slaveNum {
			atomic.StoreInt32(&s.counter, 0)
			Log.Print("=Slave Done Counter Reset: ", s.counter)
			s.runningStep <- 2
		}
		reply.State = 200
		reply.Value = nil
		return nil
	}
	return errors.New("wrong messages action.")
}
