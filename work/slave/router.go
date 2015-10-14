package slave

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

func NewRouter(net *structure.AdjacentList, c *message.LocalMessageContainer) *Router {
	r := &Router{
		net:  net,
		msgs: c,

		// processor: NewProcessor(net),
	}
	return r
}

func (r *Router) Route(msg *message.Message, reply *message.Reply) (err error) {
	switch msg.Object {
	case message.LOCAL:
		return r.toLocalAction(msg, reply)
	case message.NODE:
		return r.toNodeAction(msg, reply)
	case message.EDGE:
		return r.toEdgeAction(msg, reply)
	default:
		reply.Value, _ = json.Marshal(errors.New("object type is null"))
		reply.State = 201
		return errors.New("object type is null")
	}
}
func (r *Router) toLocalAction(msg *message.Message, reply *message.Reply) (err error) {
	// 	switch msg.Action {
	// 	case message.NEXTSTEP:

	// }
	return
}

func (r *Router) toNodeAction(msg *message.Message, reply *message.Reply) (err error) {
	switch msg.Action {
	case message.ADD:

		body := msg.UnmarshalToAddNode()
		// println(body.Id)
		n := r.net.AddNode(engine.NewObjectId(body.Id))
		if n != nil {
			reply.Value, _ = json.Marshal(n.Index.ToByte())
			reply.State = 200
			return nil
		} else {
			reply.Value = nil
			reply.State = 201
			return errors.New("text")
		}

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
	return errors.New("unknown errors.")
}

func (r *Router) toEdgeAction(msg *message.Message, reply *message.Reply) (err error) {
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
		return errors.New("text")
	}
	return nil
}
