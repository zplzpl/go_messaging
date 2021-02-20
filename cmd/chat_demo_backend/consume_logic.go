package main

import (
	"fmt"

	"go_messaging/internal/pubsub"
	"go_messaging/pkg/json"
)

type ConsumeLogic struct {
	MQPublisher *MQPublisher
}

func NewConsumeLogic(MQPublisher *MQPublisher) *ConsumeLogic {
	return &ConsumeLogic{MQPublisher: MQPublisher}
}

func (c *ConsumeLogic) ConsumeMsg(msg []byte) error {

	// parse msg
	obj := new(pubsub.Msg)
	if err := json.Unmarshal(msg, &obj); err != nil {
		return err
	}

	// handle logic
	resp := &pubsub.Msg{
		RoomId:  obj.RoomId,
		Content: []byte(fmt.Sprintf("you send msg: %s", string(obj.Content))),
	}

	buf, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// publish result to rabbitMQ
	if err := c.MQPublisher.Publish("", buf); err != nil {
		return err
	}

	return nil
}
