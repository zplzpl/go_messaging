package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go_messaging/cmd/chat_demo_backend/app/v1/model/common"
	"go_messaging/cmd/chat_demo_backend/app/v1/model/form"
	"go_messaging/pkg/unique_id"
)

type RoomHandler struct {
}

func NewRoomHandler() *RoomHandler {
	return &RoomHandler{}
}

func (t *RoomHandler) Generate(ctx *gin.Context) {

	// generate room id
	roomId := fmt.Sprintf("%d", unique_id.GenerateID().Int64())

	res := &form.RoomGenerateResp{
		RetStatus: common.RetStatus{
			Code: 0,
			Msg:  "ok",
		},
		Data: &form.RoomGenerateData{
			RoomId: roomId,
		},
	}

	ctx.JSON(http.StatusOK, res)
}
