package slave

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	// "mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/structure"
	. "mervin.me/SalonRabbit/util/log"
	"mervin.me/SalonRabbit/work"
)

import (
	"encoding/json"
	"net"
	"net/rpc"
)

var (
	MasterHost string
)

type Slave struct {
	Net *structure.AdjacentList

	MessageContainer *message.LocalMessageContainer

	DataChan   *com.DataChannel
	SignalChan *com.SignalChannel

	Router *Router

	processor *Processor

	localNet *InitLocalNet
}

func New() *Slave {
	w := new(Slave)
	msgs := message.NewLocalMessageContainer()
	//DataChannel Structure Router, call relationship
	sc := com.NewSignalChannel(w)
	dc := com.NewDataChannel(w)
	net := structure.New(dc)
	r := NewRouter(net, msgs)
	w.Net = net
	w.MessageContainer = msgs
	w.DataChan = dc
	w.SignalChan = sc

	w.Router = r
	w.processor = NewProcessor(net, sc)
	w.localNet = &InitLocalNet{r}
	println(w.Router)
	return w
}
func (s *Slave) Role() work.WorkerType {
	return work.SlaveWorker
}

func (s *Slave) Start() {
	server := rpc.NewServer()
	server.Register(s.localNet)
	server.Register(s.DataChan)
	server.Register(s.SignalChan)
	// workr.Register(new(engine.ObjectId))

	addr, err := net.ResolveTCPAddr("tcp", ":"+work.DefaultPort)
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	server.Accept(l)
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go server.ServeConn(conn)
	}
}

func (s *Slave) GetRouter() (r work.Router) {
	return s.Router
}
func (s *Slave) GetSignalController() {

}

func (s *Slave) PullSignal(msg *message.Message, reply *message.Reply) (err error) {
	Log.Print("PullSignal")
	if msg.Action == message.NEXTSTEP {
		ns := msg.UnmarshalToNextstep()
		MasterHost = ns.Host
		Log.Print(MasterHost, " ", ns.Step)
		go s.processor.Run(ns.Step, s.MessageContainer)
		// reply = &message.Reply{
		// Value: nil,
		// State: 200,
		// }
		reply.State = 200
		reply.Value = nil
	}
	if msg.Action == message.SYNC {
		//active nodes
		//vote to active
		Log.Print("SYNC REQUEST")
		s.MessageContainer.Nextstep()

		net := s.Net
		for k, _ := range s.MessageContainer.CurMsgMap {
			Log.Print(k.String())
			index, _ := k.Unmarshal()
			net.GetNodeByIndex(index).VoteToActive()
		}
		size, _ := s.MessageContainer.Size()
		// reply = &message.Reply{
		// 	Value: size,
		// 	State: 200,
		// }
		reply.State = 200
		reply.Value, _ = json.Marshal(size)
	}
	return nil
}

type InitLocalNet struct {
	Router work.Router
}

func (l *InitLocalNet) Request(msg *message.Message, reply *message.Reply) (err error) {

	err = l.Router.Route(msg, reply)
	return
}
