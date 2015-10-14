package master

import (
	"mervin.me/SalonRabbit/com/message"
	"mervin.me/SalonRabbit/engine"
	. "mervin.me/SalonRabbit/util/log"
)
import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Networking struct {
	fp           *os.File
	nodeIndexMap map[string]*engine.ObjectIndex
	workers      *WorkerManagement
}

func NewNetworking(m map[string]*engine.ObjectIndex) *Networking {
	n := &Networking{
		nodeIndexMap: m,
	}
	return n
}

func (n *Networking) Init(file string) {
	p, err := os.OpenFile(file, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	n.fp = p
	n.workers = NewWorkerManagement()
	n.workers.Start()
}

func (n *Networking) Start() {
	n.initNode()
	n.initEdge()
	n.workers.Close()
}

func (n *Networking) initNode() {
	bufIn := bufio.NewReader(n.fp)
	var delim byte = '\n'
	i := 0
	for {
		str, err := bufIn.ReadString(delim)
		if err != nil || err == io.EOF {
			break
		}
		strArr := strings.Fields(str)
		if len(strArr) == 1 || len(strArr) == 2 {
			for _, v := range strArr {
				if _, ok := n.nodeIndexMap[v]; ok {
					continue
				}
				var reply message.Reply
				msgBody := message.AddNode{
					Id: v,
				}
				args := &message.Message{
					Object: message.NODE,
					Action: message.ADD,
					Value:  msgBody.Marshal(),
				}
				for true {

					err := n.workers.Queue[i%n.workers.size].Call("InitLocalNet.Request", args, &reply)
					if err == nil && reply.State == 200 {

						var index engine.ByteIndex
						json.Unmarshal(reply.Value, &index)
						n.nodeIndexMap[v], err = index.Unmarshal()
						Log.Print("Init Node: ", reply.State)
						break
					}
				}
				i++
			} //for
			Log.Print(str)
		} //if
	} //for
	Log.Println(len(n.nodeIndexMap))
}

//
func (n *Networking) initEdge() {
	n.fp.Seek(0, 0)
	bufIn := bufio.NewReader(n.fp)
	var delim byte = '\n'
	for {
		str, err := bufIn.ReadString(delim)
		if err != nil || err == io.EOF {
			break
		}
		strArr := strings.Fields(str)
		if len(strArr) == 2 {
			src, ok1 := n.nodeIndexMap[strArr[0]]
			tar, ok2 := n.nodeIndexMap[strArr[1]]
			if ok1 && ok2 {
				var reply message.Reply
				msgBody1 := message.AddEdge{
					SrcIndex: src.ToByte(),
					TarIndex: tar.ToByte(),
				}
				args := &message.Message{
					Object: message.EDGE,
					Action: message.ADD,
					Value:  msgBody1.Marshal(),
				}
				srcIpInt := src.GetMachine()
				for true {
					err := n.workers.Get(srcIpInt).Call("InitLocalNet.Request", args, &reply)
					if err == nil && reply.State == 200 {
						Log.Print("Init Edge: ", reply.State)
						break
					}
				}

				msgBody2 := message.AddEdge{
					SrcIndex: tar.ToByte(),
					TarIndex: src.ToByte(),
				}
				args = &message.Message{
					Object: message.EDGE,
					Action: message.ADD,
					Value:  msgBody2.Marshal(),
				}
				srcIpInt = tar.GetMachine()
				for true {
					err := n.workers.Get(srcIpInt).Call("InitLocalNet.Request", args, &reply)
					if err == nil && reply.State == 200 {
						Log.Print("Init Edge: ", reply.State)
						break
					}
				}
			}

		} //if
	} //for

}
