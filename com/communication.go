package com

// import (
// 	"mervin.me/SalonRabbit/work"
// )
// import (
// 	"net"
// 	"net/rpc"
// )

// func Connect(dstAddr string) *rpc.Client {
// 	var err error
// 	for i := 0; i < 3; i++ {
// 		addr, err := net.ResolveTCPAddr("TCP", dstAddr+":"+work.DefaultPort)
// 		if err != nil {
// 			continue
// 		}
// 		conn, err := net.DialTCP("TCP", nil, addr)
// 		if err != nil {
// 			continue
// 		}
// 		c := rpc.NewClient(conn)
// 		return c
// 	}
// 	panic(err)
// }
