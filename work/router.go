package work

import (
	"mervin.me/SalonRabbit/com/message"
)

// type Router struct {
// 	nextStepMsgs *message.LocalMessageMap
// 	curStepMsgs  *message.LocalMessageMap
// 	Route        func(msg *message.Message) (reply *message.Reply)
// }

// func NewRouter(msgs [2]*message.LocalMessageMap) *Router {
// 	r := &Router{
// 		curStepMsgs:  msgs[0],
// 		nextStepMsgs: msgs[1],
// 	}
// 	return r
// }
type Router interface {
	Route(msg *message.Message, reply *message.Reply) (err error)
}
