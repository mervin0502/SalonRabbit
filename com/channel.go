package com

import (
	"mervin.me/SalonRabbit/com/message"
	. "mervin.me/SalonRabbit/util/log"
	"mervin.me/SalonRabbit/work"
	// "mervin.me/SalonRabbit/work/slave"
	// "mervin.me/util/ip"
)

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strconv"
)

type Channel struct {
	// conn *rpc.Client
	Worker work.Worker
}

type DataChannel struct {
	Channel
	// Router work.Router
}
type SignalChannel struct {
	Channel
	// Signal chan<- *message.Message
}

//
//
func NewDataChannel(w work.Worker) *DataChannel {
	Log.Print("New DataChannel")
	c := &DataChannel{
		Channel: Channel{
			Worker: w,
		},
		// Router: w.GetRouter(),
	}

	return c
}

func (d *DataChannel) Send(msg *message.Message, reply *message.Reply) (err error) {
	Log.Print("DataChannel.Send")
	host := getTargetHost(string(msg.Target[0:8]))
	conn := connectTo(host)
	if conn == nil {
		return err
	}
	err = conn.Call("DataChannel.Request", msg, reply)
	conn.Close()
	return err
}
func (d *DataChannel) Request(msg *message.Message, reply *message.Reply) (err error) {
	Log.Print("DataChannel.Request")
	// println(d.Router)
	err = d.Worker.GetRouter().Route(msg, reply)
	return
}

///
///
///
///

func NewSignalChannel(w work.Worker) *SignalChannel {
	Log.Print("New SignalChannel")
	c := &SignalChannel{
		Channel: Channel{
			Worker: w,
		},
	}
	return c
}

func (s *SignalChannel) Send(msg *message.Message, reply *message.Reply) (err error) {
	if len(msg.Target) < 8 {
		panic(errors.New("msg.Target error."))
	}
	host := getTargetHost(string(msg.Target[0:8]))
	conn := connectTo(host)
	if conn == nil {
		return errors.New("connect error.")
	}

	Log.Print("SignalChannel Send: ", host)
	// var r message.Reply
	err = conn.Call("SignalChannel.Request", msg, reply)
	Log.Print(reply.State)
	// *reply = r
	conn.Close()
	return err
}
func (s *SignalChannel) Request(msg *message.Message, reply *message.Reply) (err error) {
	Log.Print("SignalChannel Request")
	if msg.Object == message.LOCAL && msg.Action == message.NEXTSTEP && s.Worker.Role() == work.SlaveWorker {
		// master->slave
		err = s.Worker.PullSignal(msg, reply)
		Log.Print("NEXTSTEP:Master->Slave ", reply.State)

	} else if msg.Object == message.LOCAL && msg.Action == message.SYNC && s.Worker.Role() == work.SlaveWorker {
		//master->slave
		//sync
		err = s.Worker.PullSignal(msg, reply)
		Log.Print("SYNC:Master->Slave ", reply.State)

	} else if msg.Object == message.NETWWORK && msg.Action == message.SLAVEDONE && s.Worker.Role() == work.MasterWorker {
		//slave->master
		//
		err = s.Worker.PullSignal(msg, reply)
		Log.Print("SLAVEDONE:Slave->Master ", reply.State)

	} else {
		return errors.New("wrong action.")
	}
	Log.Print("Requst Done ", reply.State)
	return
}
func getTargetHost(hex string) string {
	Log.Print("getTargetHost")
	v := hex[0:2]
	a, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		panic(err)
	}
	v = hex[2:4]
	b, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		panic(err)
	}
	v = hex[4:6]
	_c, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		panic(err)
	}
	v = hex[6:8]
	d, err := strconv.ParseInt(v, 16, 32)
	if err != nil {
		panic(err)
	}
	host := fmt.Sprintf("%d.%d.%d.%d", a, b, _c, d)
	return host
}
func connectTo(dstAddr string) *rpc.Client {
	Log.Print("connectTo: " + dstAddr)
	for i := 0; i < 3; i++ {
		addr, err := net.ResolveTCPAddr("tcp", dstAddr+":"+work.DefaultPort)
		if err != nil {
			panic(err)
			continue
		}
		con, err := net.DialTCP("tcp", nil, addr)
		if err != nil {
			panic(err)
			continue
		}
		conn := rpc.NewClient(con)

		return conn
	}
	return nil
}
