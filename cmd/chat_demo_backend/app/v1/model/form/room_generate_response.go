package form

import "go_messaging/cmd/chat_demo_backend/app/v1/model/common"

type RoomGenerateResp struct {
	common.RetStatus
	Data *RoomGenerateData `json:"data"`
}

type RoomGenerateData struct {
	RoomId string `json:"room_id"`
}
