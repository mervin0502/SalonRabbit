package master

import (
	// "mervin.me/SalonRabbit/engine/message"
	"mervin.me/SalonRabbit/work"
	"mervin.me/util/ip"
)

import (
	"net"
	"net/rpc"
	// "strconv"
	"sync"
)

var SlaveHosts = []string{
	"192.168.36.11",
	"192.168.36.12",
	// "192.168.36.13",
}

type WorkerManagement struct {
	Queue     []*rpc.Client
	ip2RPCMap map[uint32]int
	size      int
	mutex     *sync.Mutex
}

func NewWorkerManagement() *WorkerManagement {
	s := len(SlaveHosts)
	wm := &WorkerManagement{
		Queue:     make([]*rpc.Client, s),
		ip2RPCMap: make(map[uint32]int, s),
		size:      s,
		mutex:     new(sync.Mutex),
	}
	return wm
}

func (w *WorkerManagement) Start() {
	for i, workerHost := range SlaveHosts {
		addr, err := net.ResolveTCPAddr("tcp", workerHost+":"+work.DefaultPort)
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", nil, addr)
		// defer conn.Close()

		c := rpc.NewClient(conn)
		// defer c.Close()
		w.Queue[i] = c
		w.ip2RPCMap[ip.Ip2Int(workerHost)] = i
		// println(ip.Ip2Int(workerHost), i)
	}
}

func (w *WorkerManagement) Get(ipInt uint32) *rpc.Client {
	return w.Queue[w.ip2RPCMap[ipInt]]
}
func (w WorkerManagement) Size() int {
	return len(w.Queue)
}
func (w *WorkerManagement) Close() {
	for _, worker := range w.Queue {
		worker.Close()
	}
}
