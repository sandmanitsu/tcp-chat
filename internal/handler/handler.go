package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tcp-chat/internal/broadcast"
	"tcp-chat/internal/client"
	"tcp-chat/internal/room"
)

func HandleConn(conn net.Conn, cast broadcast.Broadcast) {
	ch := make(chan string)
	addr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)
	go writer(conn, ch)

	// creating user
	var name string
	fmt.Fprintf(conn, "Enter your name:")
	for input.Scan() {
		name = input.Text()
		if name == "" {
			ch <- "Your name is empty! Try again..."
			continue
		}

		fmt.Printf("user %s name: %s entered to chat\n", addr, name)
		break
	}
	client := client.Client{
		Name:     name,
		CurrRoom: "",
	}

	// connect to main room
	cast.Entering <- broadcast.EnteringData{
		Chan:     ch,
		Client:   &client,
		RoomName: room.MainRoomName,
	}

	// listen msg from user
	var isCreating bool
	for input.Scan() {
		msg := input.Text()

		if isCreating {
			cast.Creating <- broadcast.CreatingData{
				Chan:     ch,
				Client:   &client,
				RoomName: msg,
			}
			isCreating = false

			continue
		}

		switch {
		case strings.Contains(msg, room.Help):
			ch <- room.GetCommandsDescription(cast.Room[client.CurrRoom].Commands)
			continue
		case strings.Contains(msg, room.GetUsers):
			ch <- room.GetUsersInRoom(cast.Room[client.CurrRoom])
			continue
		case strings.Contains(msg, room.CreateRoom):
			ch <- "Enter room name:"
			isCreating = true
			continue
		}

		cast.Messages <- broadcast.Message{
			Msg:      fmt.Sprintf("%s: %s", name, msg),
			RoomName: client.CurrRoom,
		}
	}

	// user leaving
	cast.Leaving <- broadcast.LeavingData{
		Chan:   ch,
		Client: &client,
		IsExit: true,
	}
	// cast.Messages <- fmt.Sprintf("%s exit the chat!\n", name)
	conn.Close()
}

func writer(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, "%s\n", msg)
	}
}
