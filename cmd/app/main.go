package main

import (
	"log"
	"net"
	"tcp-chat/internal/broadcast"
	"tcp-chat/internal/handler"
	"tcp-chat/internal/room"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Chat started")
	broadcast := broadcast.NewMainBroardcast(room.NewMainRooms())
	go broadcast.Broadcast()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handler.MainHandleConn(conn, broadcast)
	}
}
