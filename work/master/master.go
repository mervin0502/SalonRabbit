package master

//init network
//opration on network
//node/edge
//put delete update get

/*
  init network from text or others
  ping the worker, recovery
  run the Compute function of node
  aggregator
*/

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	. "mervin.me/SalonRabbit/util/log"
	"mervin.me/SalonRabbit/work"
)

import (
	"net"
	"net/rpc"
	"strings"
)

var (
	LocalIP string
)

func init() {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	LocalIP = strings.Split(conn.LocalAddr().String(), ":")[0]
}

type Master struct {
	NodeIndexMap map[string]*engine.ObjectIndex

	MessageContainer *message.LocalMessageContainer
	DataChan         *com.DataChannel
	SignalChan       *com.SignalChannel
	superstep        *Superstep
}

//New
func New() *Master {
	w := new(Master)
	msgs := message.NewLocalMessageContainer()
	sc := com.NewSignalChannel(w)

	w.NodeIndexMap = make(map[string]*engine.ObjectIndex, 100)
	w.MessageContainer = msgs
	w.DataChan = com.NewDataChannel(w)
	w.SignalChan = sc
	w.superstep = NewSuperstep(sc)
	return w
}
func (m *Master) Role() work.WorkerType {
	return work.MasterWorker
}

//init network
func (m *Master) InitNetwork(file string) {
	net := NewNetworking(m.NodeIndexMap)
	net.Init(file)
	net.Start()
	Log.Println("InitNetwork Done")
}
func (m *Master) Start() {
	//
	go m.superstep.Run()
	//
	server := rpc.NewServer()
	server.Register(m.DataChan)
	server.Register(m.SignalChan)
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

func (m *Master) GetSignalController() {

}
func (m *Master) GetRouter() work.Router {
	return nil
}
func (m *Master) PullSignal(msg *message.Message, reply *message.Reply) (err error) {
	return m.superstep.Put(msg, reply)
}
