package main

import (
	"go_messaging/internal/pubsub"
	"go_messaging/internal/ws_connect"
	"go_messaging/pkg/json"
)

type ConsumeLogic struct {
	RoomManager *ws_connect.RoomManager
}

func NewConsumeLogic(roomManager *ws_connect.RoomManager) *ConsumeLogic {
	return &ConsumeLogic{RoomManager: roomManager}
}

func (c *ConsumeLogic) ConsumeMsg(msg []byte) error {

	// parse msg
	obj := new(pubsub.Msg)
	if err := json.Unmarshal(msg, &obj); err != nil {
		return err
	}

	room, found := c.RoomManager.PureFindRoom(obj.RoomId)
	if !found {
		return nil
	}

	// send msg to room
	room.Broadcast(obj.Content)

	return nil
}
