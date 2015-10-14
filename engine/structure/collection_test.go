package structure

import (
	"testing"
)

import (
	"mervin.me/SalonRabbit/engine"
)

func TestNew(t *testing.T) {
	col := NewCollection()
	n, err := col.Add(engine.NewObjectId("id1"))
	t.Log(err)
	t.Log(n.Id.Value)
	t.Log(n.Index.String())
}
