package edge

import (
	"mervin.me/SalonRabbit/engine"
	"mervin.me/SalonRabbit/engine/attribute"
)

type Edge struct {
	// Id     engine.ObjectId
	Source *engine.ObjectIndex
	Target *engine.ObjectIndex

	Attributes *attribute.Attribute
}

//New return the object of struct Edge
func New(src, tar *engine.ObjectIndex) *Edge {
	e := &Edge{
		Source: src,
		Target: tar,
	}
	return e
}

//GetTargetNode returns the target node of edge e
func (e *Edge) GetTargetNode() *engine.ObjectIndex {
	if e == nil {
		return nil
	}
	return e.Target
}

//
func (e *Edge) PutAttribte(key string, value interface{}) {
	e.Attributes.Put(key, value)
}
func (e *Edge) GetAttribute(key string) interface{} {
	return e.Attributes.Get(key)
}
