package master

import (
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/structure"
)
import (
	"encoding/json"
	"errors"
)

type Router struct {
	net  *structure.AdjacentList
	msgs *message.LocalMessageContainer
}

func NewRouter(net *structure.AdjacentList, msgs *message.LocalMessageContainer) *Router {
	r := &Router{
		net: net,
		// curStepMsgs:  *msgs[0],
		msgs: msgs,

		// processor: NewProcessor(net),
	}
	return r
}

func (r *Router) Route(msg *message.Message) (reply *message.Reply) {
	switch msg.Object {
	case message.LOCAL:
		return r.toLocalAction(msg)
	case message.NODE:
		return r.toNodeAction(msg)
	case message.EDGE:
		return r.toEdgeAction(msg)
	default:
		reply.Value, _ = json.Marshal(errors.New("object type is null"))
		reply.State = 201
	}
	return
}
func (r *Router) toLocalAction(msg *message.Message) (reply *message.Reply) {
	// 	switch msg.Action {
	// 	case message.NEXTSTEP:

	// }
	return
}

func (r *Router) toNodeAction(msg *message.Message) (reply *message.Reply) {
	switch msg.Action {
	case message.ADD:
		body := msg.UnmarshalToAddNode()
		r.net.AddNode(engine.NewObjectId(body.Id))
	case message.DEL:
		r.msgs.Put(msg.Target, msg)
		reply.Value = nil
		reply.State = 200
		return nil
	case message.PUT:
		r.msgs.Put(msg.Target, msg)
		reply.Value = nil
		reply.State = 200
		return nil
	case message.GET:
		r.msgs.Put(msg.Target, msg)
		reply.Value = nil
		reply.State = 200
		return nil
	default:
		v, _ := json.Marshal(errors.New("action type is null"))
		reply = &message.Reply{
			Value: v,
			State: 201,
		}
		return
	}
	reply = &message.Reply{
		Value: nil,
		State: 200,
	}
	return reply
}

func (r *Router) toEdgeAction(msg *message.Message) (reply *message.Reply) {
	switch msg.Action {
	case message.ADD:
		// index, _ := engine.UnmarshalObjectIndex(msg.Target)
		// r.nextStepMsgs.Put(index, msg)
		body := msg.UnmarshalToAddEdge()
		r.net.AddEdge(body)
		reply.Value = nil
		reply.State = 200
		return nil
	case message.DEL:
		r.msgs.Put(msg.Target, msg)

		reply.Value = nil
		reply.State = 200
		return nil
	case message.PUT:
		r.msgs.Put(msg.Target, msg)

		reply.Value = nil
		reply.State = 200
		return nil
	case message.GET:
		r.msgs.Put(msg.Target, msg)

		reply.Value = nil
		reply.State = 200
		return nil
	default:
		reply.Value, _ = json.Marshal(errors.New("object type is null"))
		reply.State = 201
		return
	}
	reply = &message.Reply{
		Value: nil,
		State: 200,
	}
	return reply
}
