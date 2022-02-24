package main

import (
	"gin_chat/server"
	"log"
)

func main() {
	s := server.NewSerer()

	log.Fatal(s.Run())
}

func InitConfig() {
	
}
