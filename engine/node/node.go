package node

import (
	"mervin.me/SalonRabbit/com"
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/attribute"
	"mervin.me/SalonRabbit/engine/edge"
)

import (
	// "fmt"
	// "strconv"
	"sync"
)

type NodeState uint8

const (
	ACTIVE NodeState = iota
	INACTIVE
)

var (
	MessageChan = make(chan *message.Message, 100)
)

//
// receive the superstep message from the master
// traverse the local node set
// if node's status is active
// execute the Compute() function

// vote itself to halt
// new the message and push to the message list of local workr

// combiner the messages
// send messages body
//
// receive messages and insert to receive-messages list
// decompose the messages and apply every message to the destination node
// set the node's status is active if it receive message
//
// send message to master and done the current superstep work
// tell master how many node is active

//
// Compute():
//		get messages
//    decompose the messages and apply msg to node or edge
//		create newly message
//		send to message
//		vote to halt
//
type Node struct {
	Index     *engine.ObjectIndex
	Id        *engine.ObjectId
	Superstep uint32

	Attributes *attribute.Attribute
	State      NodeState

	OutEdges []*edge.Edge

	OutDegree int

	mutex *sync.Mutex

	dataChan *com.DataChannel

	Compute func(msgs *message.MessageCollection)
}

func New(id *engine.ObjectId) *Node {
	n := &Node{
		Id:         id,
		Attributes: attribute.New(),
		State:      ACTIVE,
		OutEdges:   make([]*edge.Edge, 0),
		mutex:      new(sync.Mutex),
	}

	return n
}

func (n *Node) SetDataChan(dc *com.DataChannel) {
	n.dataChan = dc
}

func (n *Node) SendMessageTo(msg *message.Message, reply *message.Reply) (err error) {

	return n.dataChan.Send(msg, reply)
}

//VoteToHalt
func (n *Node) VoteToHalt() {
	if n.State == ACTIVE {
		n.State = INACTIVE
	}
}
func (n *Node) VoteToActive() {
	if n.State == INACTIVE {
		n.State = ACTIVE
	}
}

//AddEdge
func (n *Node) AddEdge(tarNodeIndex *engine.ObjectIndex) (*edge.Edge, bool) {
	if n.OutDegree > 0 {
		for _, e := range n.OutEdges {
			if e.Target.Equals(tarNodeIndex) {
				// println("Node.AddEdge.False")
				return e, false
			}
		}
	}
	e := edge.New(n.Index, tarNodeIndex)
	println(n.Index.String(), tarNodeIndex.String())
	n.OutEdges = append(n.OutEdges, e)
	n.OutDegree++
	return e, true
}
func (n *Node) GetEdge(tarNodeIndex *engine.ObjectIndex) (*edge.Edge, bool) {
	if n.OutDegree > 0 {
		for _, e := range n.OutEdges {
			if e.Target.Equals(tarNodeIndex) {
				// println("Node.AddEdge.False")
				return e, true
			}
		}
	}
	return nil, false
}

func (n *Node) DelEdge(tarNodeIndex *engine.ObjectIndex) *edge.Edge {
	if n.OutDegree > 0 {
		for k, e := range n.OutEdges {
			if e.Target.Equals(tarNodeIndex) {
				n.OutEdges = append(n.OutEdges[:k], n.OutEdges[k+1:]...)
				return e
			}
		}
	}
	return nil
}
func (n *Node) IterOutEdge() <-chan *edge.Edge {
	e := make(chan *edge.Edge, 10)
	go func(edge chan<- *edge.Edge) {
		for _, v := range n.OutEdges {
			edge <- v
		}
		close(edge)
	}(e)
	return e
}

//
func (n *Node) PutAttribte(key string, value interface{}) {
	n.Attributes.Put(key, value)
}
func (n *Node) GetAttribute(key string) interface{} {
	return n.Attributes.Get(key)
}
