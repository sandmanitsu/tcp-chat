package handler

import (
	"bufio"
	"fmt"
	"net"
	"slices"
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
	isCreating := false
	isEntering := false
	for input.Scan() {
		msg := input.Text()

		if strings.Contains(msg, room.Help) {
			ch <- room.GetCommandsDescription(cast.Room[client.CurrRoom].Commands)
		}

		if isCreating {
			if msg == "" {
				ch <- "Empty room name. Try again!"

				isCreating = false

				continue
			}

			cast.Creating <- broadcast.CreatingData{
				Chan:     ch,
				Client:   &client,
				RoomName: msg,
			}
			isCreating = false

			continue
		}

		if isEntering {
			_, ok := cast.Room[msg]
			if msg == "" || !ok {
				ch <- "Empty room name. Try again!"

				isCreating = false

				continue
			}

			cast.Entering <- broadcast.EnteringData{
				Chan:     ch,
				Client:   &client,
				RoomName: msg,
			}

			isEntering = false

			continue
		}

		if slices.Contains(cast.Room[client.CurrRoom].Commands, msg) {
			switch {
			case strings.Contains(msg, room.GetUsers):
				ch <- room.GetUsersInRoom(cast.Room[client.CurrRoom])
			case strings.Contains(msg, room.CreateRoom):
				ch <- "Enter room name:"
				isCreating = true
			case strings.Contains(msg, room.EnterChat):
				ch <- fmt.Sprintf("Choose room:\n%s", room.GetRoomsNames(cast.Room))
				isEntering = true
			case strings.Contains(msg, room.LeaveChat):
				cast.Leaving <- broadcast.LeavingData{
					Chan:   ch,
					Client: &client,
					IsExit: false,
				}
			}
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

	cast.Messages <- broadcast.Message{
		Msg:      fmt.Sprintf("%s exit the chat!\n", name),
		RoomName: client.CurrRoom,
	}

	conn.Close()
}

func writer(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, "%s\n", msg)
	}
}
