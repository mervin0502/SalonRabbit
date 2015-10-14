package message

import (
	"mervin.me/SalonRabbit/engine"
	. "mervin.me/SalonRabbit/util/log"
)

type entry struct {
	msg  *Message
	next *entry
}
type MessageQueue struct {
	front *entry
	rear  *entry
	size  int
}

func (q *MessageQueue) Push(msg *Message) {
	Log.Print("MessageQueue.Push")
	e := &entry{
		msg:  msg,
		next: nil,
	}
	if q.front == nil && q.rear == nil {
		q.front = e
		q.rear = e
	} else {
		q.rear.next = e
		q.rear = e
	}
	q.size++
}
func (q *MessageQueue) Pop() *Message {
	Log.Print("Pop")
	if q.front != nil {
		var msg *Message
		if q.front == q.rear {
			//one entry
			msg = q.front.msg
			q.front = nil
			q.rear = nil
		} else {
			msg = q.front.msg
			temp := q.front.next
			q.front.next = nil
			q.front = temp
		}

		q.size--
		return msg
	}
	return nil

}
func (q *MessageQueue) HasNext() bool {
	Log.Print("HasNext")
	return q.front != nil
}
func (q *MessageQueue) Size() int {
	Log.Print("Size")
	return q.size
}

type MessageCollection struct {
	NodeMsg *MessageQueue
	EdgeMsg *MessageQueue
}

// type LocalMessageMap map[*engine.ByteIndex]*MessageCollection
type LocalMessageContainer struct {
	CurMsgMap  map[engine.ByteIndex]*MessageCollection
	NextMsgMap map[engine.ByteIndex]*MessageCollection
}

func NewLocalMessageContainer() *LocalMessageContainer {
	c := &LocalMessageContainer{
		CurMsgMap:  make(map[engine.ByteIndex]*MessageCollection, 0),
		NextMsgMap: make(map[engine.ByteIndex]*MessageCollection, 0),
	}
	return c
}

// func NewLocalMessageMap() *LocalMessageMap {
// 	mm := make(map[*engine.ByteIndex]*MessageCollection, 100)
// 	mm2 := LocalMessageMap(mm)
// 	return &mm2
// }

func (l *LocalMessageContainer) Put(index engine.ByteIndex, msg *Message) {
	Log.Print("LocalMessageContainer.Put")
	c := l.NextMsgMap
	if msg.Object == NODE {
		if v, ok := c[index]; ok {
			v.NodeMsg.Push(msg)
		} else {
			v := &MessageQueue{
				front: nil,
				rear:  nil,
				size:  0,
			}
			v.Push(msg)
			v1 := new(MessageCollection)
			v1.NodeMsg = v
			c[index] = v1
		}
	} else if msg.Object == EDGE {
		if v, ok := c[index]; ok {
			v.EdgeMsg.Push(msg)
		} else {
			v := new(MessageQueue)
			v.Push(msg)
			v1 := new(MessageCollection)
			v1.EdgeMsg = v
			c[index] = v1
		}
	} else {
	}

}

func (l LocalMessageContainer) Get(index engine.ByteIndex) *MessageCollection {
	Log.Print("Get")
	c := l.CurMsgMap
	if v, ok := c[index]; ok {
		return v
	} else {
		return nil
	}
}

func (l *LocalMessageContainer) Nextstep() {
	Log.Print("Nextstep")
	temp := make(map[engine.ByteIndex]*MessageCollection, 0)
	l.CurMsgMap = l.NextMsgMap
	l.NextMsgMap = temp
}

func (l *LocalMessageContainer) Size() (curMsgSize int, nextMsgSize int) {
	Log.Print("Size")
	curMsgSize = len(l.CurMsgMap)
	nextMsgSize = len(l.NextMsgMap)
	return
}
