package structure

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/edge"
	"mervin.me/SalonRabbit/engine/node"
)

import (
// "sync"
)

type AdjacentList struct {
	CurSuperStep uint32
	nodeNum      uint32
	edgeNum      uint32
	// mutex        *sync.Mutex
	// rwmutex      *sync.RWMutex
	dataChan *com.DataChannel
	topo     *Collection
}

//New one newly adjacent-list
func New(dc *com.DataChannel) *AdjacentList {
	a := &AdjacentList{
		nodeNum: 0,
		edgeNum: 0,
		// mutex:   new(sync.Mutex),
		// rwmutex: new(sync.RWMutex),
		dataChan: dc,
		topo:     NewCollection(),
	}
	return a
}

//GetNode
func (a *AdjacentList) GetNodeById(id *engine.ObjectId) *node.Node {
	n, _ := a.topo.GetById(id)
	return n
}
func (a *AdjacentList) GetNodeByIndex(index *engine.ObjectIndex) *node.Node {
	n, _ := a.topo.Get(index)
	return n
}

//
func (a *AdjacentList) GetCollection() *Collection {
	return a.topo
}

//AddNode
func (a *AdjacentList) AddNode(id *engine.ObjectId) *node.Node {
	// a.mutex.Lock()
	// defer a.mutex.Unlock()
	//println(id.Value)
	n, ok := a.topo.Add(id)
	if ok {
		n.Superstep = a.CurSuperStep
		n.SetDataChan(a.dataChan)
	}
	a.nodeNum += 1

	return n
}

//AddEdge
func (a *AdjacentList) AddEdge(msg *message.AddEdge) (*edge.Edge, bool) {
	srcIndex, _ := msg.SrcIndex.Unmarshal()
	tarIndex, _ := msg.TarIndex.Unmarshal() // println("Func:AdjacentList.AddEdge")
	srcNode, _ := a.topo.Get(srcIndex)
	// println(srcNode.Index.GetMachine())
	e, ok := srcNode.AddEdge(tarIndex)
	return e, ok
}
