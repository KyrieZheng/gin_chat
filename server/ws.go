package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ws struct {
	upgrader *websocket.Upgrader
}

var Websocket = &ws{
	upgrader: &websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	},
}

func (s *ws) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")
		conn, err := s.upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}

		Room.MsgJoin(name)
		control := Room.Join(name)
		defer control.Leave()

		newMessage := make(chan string)
		go func ()  {
			var res = struct {Msg string `json:"msg"`} {}
			for {
				err := conn.ReadJSON(&res)
				if err != nil {
					close(newMessage)
					return
				}
				newMessage <- res.Msg
			}
		}()

		for {
			select {
			case event := <-control.Pipe:
				if conn.WriteJSON(&event) != nil {
					return
				}
			case msg, ok := <-newMessage:
				if !ok {
					return
				}
				fmt.Println("revice:", msg)
				control.Say(msg)
			}
		}
	}
}