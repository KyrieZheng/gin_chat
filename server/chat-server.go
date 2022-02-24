package server

import (
	"gin_chat/core"

	"github.com/gin-gonic/gin"
)


var Room = core.NewRoom()

func NewSerer() *gin.Engine {
	s := gin.Default()

	// static files
	s.Static("/static", "./static")
	s.StaticFile("/", "web/index.html")
	s.StaticFile("/refresh", "./web/refresh.html")
	s.StaticFile("/polling", "./web/polling.html")
	s.StaticFile("/ws", "./web/ws.html")

	s.GET("/ws/socket", Websocket.Handle())

	return s
}