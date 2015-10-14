package slave

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/node"
	"mervin.me/SalonRabbit/engine/structure"
	. "mervin.me/SalonRabbit/util/log"
	"mervin.me/util/ip"
)

import (
	// "sync"
	"encoding/json"
)

type Processor struct {
	net        *structure.AdjacentList
	signalChan *com.SignalChannel
}

func NewProcessor(net *structure.AdjacentList, signal *com.SignalChannel) *Processor {
	Log.Print("NEW")
	p := &Processor{
		net:        net,
		signalChan: signal,
	}
	return p
}
func (p *Processor) Run(superstep uint32, msgs *message.LocalMessageContainer) {
	Log.Print("RUN")
	net := p.net
	net.CurSuperStep = superstep

	c := net.GetCollection()
	i := c.Iter()
	for {
		n, ok := <-i
		if !ok {
			break
		}
		if n.State == node.ACTIVE {
			n.Superstep = superstep
			n.Compute = func(msg *message.MessageCollection) {
				n.Superstep = superstep
				key := "max"
				flag := false
				var max int
				if n.Superstep == 1 {
					flag = true
					max = (n.OutDegree)
					n.PutAttribte(key, max)
				} else {
					nmsg := msg.NodeMsg
					v1 := n.GetAttribute(key).(int)
					max = v1
					for nmsg.HasNext() {
						m := nmsg.Pop()
						c := m.UnmarshalToPutNodeAttr()
						var v2 int
						json.Unmarshal(c.Value, &v2)

						if max < v2 {
							flag = true
							max = v2
						}
					} //for

				}

				if flag {
					n.PutAttribte(key, max)
					Log.Print("Max: ", n.Id.Value, "\t", max)
					e := n.IterOutEdge()
					for {
						edge, ok := <-e
						if !ok {
							break
						}
						v, _ := json.Marshal(max)
						body := &message.PutNodeAttr{
							Index: edge.GetTargetNode().ToByte(),
							Key:   key,
							Value: v,
						}
						msg := &message.Message{
							Object: message.NODE,
							Action: message.PUT,
							Target: edge.GetTargetNode().ToByte(),
							Value:  body.Marshal(),
						}
						var reply message.Reply
						for {
							err := n.SendMessageTo(msg, &reply)
							if err == nil && reply.State == 200 {
								break
							}
						} //for

					} //for iter edge
				}
				n.VoteToHalt()
			} //compute

			n.Compute(msgs.Get(n.Index.ToByte()))
		}
	}
	//all node done
	p.sendSlaveDone()
}

func (p *Processor) sendSlaveDone() {
	Log.Print("sendSlaveDone")
	var target engine.ByteIndex
	copy(target[0:8], ip.Ip2Hex(MasterHost))
	msg := &message.Message{
		Object: message.NETWWORK,
		Action: message.SLAVEDONE,
		Target: target,
		Value:  nil,
	}
	var reply message.Reply
	for {
		err := p.signalChan.Send(msg, &reply)
		if reply.State == 200 && err == nil {
			break
		} else if err != nil {
			println(err.Error())
		} else {
		}
	}
}
