package work

import (
	"mervin.me/SalonRabbit/com/message"
)

type WorkerType int8

const (
	MasterWorker WorkerType = iota
	SlaveWorker
)

type Worker interface {
	Role() WorkerType // master or slave
	// SetRouter(r Router)
	GetRouter() (r Router)
	PullSignal(msg *message.Message, reply *message.Reply) (err error)
	GetSignalController()
}
