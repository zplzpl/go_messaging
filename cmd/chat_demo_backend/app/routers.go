package api

import (
	"github.com/gin-gonic/gin"
	"go_messaging/cmd/chat_demo_backend/app/v1/handlers"
)

func RegisterRouter(router *gin.Engine) {

	v1Group := router.Group("/api/v1")

	roomGroup := v1Group.Group("/room")
	{
		room := handlers.NewRoomHandler()
		roomGroup.POST("/generate", room.Generate)
	}

}
