package engine

import ()

//Type ObjectId
type ObjectId struct {
	Value string
	hash  uint32
}

//NewObjectId
func NewObjectId(str string) *ObjectId {
	o := &ObjectId{
		Value: str,
	}
	o.hash = 17
	for _, i := range str {
		o.hash += o.hash*23 + uint32(i)
	}
	return o
}

//Hash
func (o *ObjectId) Hash() uint32 {
	return o.hash
}

//Equals
func (o *ObjectId) Equals(i *ObjectId) bool {
	return i.Value == o.Value
}
