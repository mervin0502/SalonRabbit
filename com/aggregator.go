package com

import (
	"mervin.me/SalonRabbit/com/message"
)

type Aggregator struct {
	dataChan  *DataChannel
	Value     *message.Message
	Aggregate func(msgs *message.MessageCollection) interface{}
}

func NewAggregator(dc *DataChannel) *Aggregator {
	a := &Aggregator{
		dataChan: dc,
	}
	return a
}
func (a *Aggregator) Min() {

}

func (a *Aggregator) Max() {

}

func (a *Aggregator) Average() {

}

func (a *Aggregator) Sum() {

}
func (a *Aggregator) Size() {

}

//Send sends a message to other slave or master
func (a *Aggregator) SendMessageTo() {

}
