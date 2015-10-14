package message

import (
	"mervin.me/SalonRabbit/engine"
	// . "mervin.me/SalonRabbit/util/log"
)

import (
	"encoding/json"
)

//
//
//  josn.rwamessage
//
//

type ActionType uint8

const (
	ADD ActionType = iota //add a node/edge
	DEL                   //del a node/edge/attribute
	PUT                   //put the attribute value of node/edge
	GET                   //get the attribute value of node/edge
	NEXTSTEP
	SLAVEDONE
	SYNC
)

type ObjectType uint8

const (
	NETWWORK ObjectType = iota
	LOCAL
	NODE
	EDGE
)

type Message struct {
	Object ObjectType       //operation object
	Action ActionType       //action of object
	Target engine.ByteIndex //the destination node index
	Value  json.RawMessage
}

type Reply struct {
	Value json.RawMessage
	// Addr  uint32
	State uint8
}

type Attribute struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

//ADD NODE
type AddNode struct {
	Id         string       `json:"id"`
	Attributes []*Attribute `json:"attributes"`
}

//DEL NODE
type DelNode struct {
	Id    string           `json:"id"`
	Index engine.ByteIndex `json:"index"`
}

//PUT NODE
type PutNodeAttr struct {
	Index engine.ByteIndex `json:"index"`
	Key   string           `json:"key"`
	Value json.RawMessage  `json:"value"`
	// Value int `json:"value"`
}

//GET NODE
type GetNodeAttr struct {
	Index engine.ByteIndex `json:"index"`
	Key   string           `json:"key"`
}

//
//ADD EDGE
type AddEdge struct {
	SrcIndex   engine.ByteIndex `json:"src"`
	TarIndex   engine.ByteIndex `json:"tar"`
	Attributes []*Attribute     `json:"attributes"`
}

//DEL EDGE
type DelEdge struct {
	SrcIndex engine.ByteIndex `json:"src"`
	TarIndex engine.ByteIndex `json:"tar"`
}

//PUT EDGE
type PutEdgeAttr struct {
	SrcIndex engine.ByteIndex `json:"src"`
	TarIndex engine.ByteIndex `json:"tar"`
	Key      string           `json:"key"`
	Value    json.RawMessage  `json:"value"`
}

//GET EDGE
type GetEdgeAttr struct {
	SrcIndex engine.ByteIndex `json:"src"`
	TarIndex engine.ByteIndex `json:"tar"`
	Key      string           `json:"key"`
}

//
//
type Nextstep struct {
	Host string `json:"host"`
	Step uint32 `json:"step"`
}

//
//
//
func (p *AddNode) Marshal() []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (p *DelNode) Marshal() []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}
func (p *PutNodeAttr) Marshal() []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}
func (p *GetNodeAttr) Marshal() []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}
func (p *AddEdge) Marshal() []byte {
	bytes, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (n *Nextstep) Marshal() []byte {
	bytes, err := json.Marshal(n)
	if err != nil {
		panic(err)
	}
	return bytes
}

//
//
//
//

func (m *Message) UnmarshalToAddNode() *AddNode {
	var body AddNode
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToDelNode() *DelNode {
	var body DelNode
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToPutNodeAttr() *PutNodeAttr {
	var body PutNodeAttr
	json.Unmarshal(m.Value, &body)

	return &body
}
func (m *Message) UnmarshalToGetNodeAttr() *GetNodeAttr {
	var body GetNodeAttr
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToAddEdge() *AddEdge {
	var body AddEdge
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToDelEdge() *DelEdge {
	var body DelEdge
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToPutEdgeAttr() *PutEdgeAttr {
	var body PutEdgeAttr
	json.Unmarshal(m.Value, &body)
	return &body
}
func (m *Message) UnmarshalToGetEdgeAttr() *GetEdgeAttr {
	var body GetEdgeAttr
	json.Unmarshal(m.Value, &body)
	return &body
}

func (m *Message) UnmarshalToNextstep() *Nextstep {
	var body Nextstep
	json.Unmarshal(m.Value, &body)
	return &body
}
