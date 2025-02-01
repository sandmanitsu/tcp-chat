package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tcp-chat/internal/broadcast"
	"tcp-chat/internal/client"
)

func MainHandleConn(conn net.Conn, cast *broadcast.MainBroadcast) {
	ch := make(chan string)
	addr := conn.RemoteAddr().String()
	input := bufio.NewScanner(conn)
	go writer(conn, ch)

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
		Name: name,
	}

	ch <- fmt.Sprintf("Welcome, %s!\n", name)
	cast.Entering <- broadcast.EnteringData{
		Chan:   ch,
		Client: client,
	}
	cast.Messages <- fmt.Sprintf("%s enter the chat!\n", name)

	var isCreating bool
	for input.Scan() {
		msg := input.Text()

		if isCreating {
			cast.Creating <- broadcast.CreatingData{
				Chan:     ch,
				Client:   client,
				RoomName: msg,
			}
			isCreating = false

			continue
		}

		switch {
		case strings.Contains(msg, broadcast.Help):
			ch <- broadcast.GetCommandsDescription(cast.Commands)
			continue
		case strings.Contains(msg, broadcast.GetUsers):
			ch <- cast.GetUsers()
			continue
		case strings.Contains(msg, broadcast.CreateRoom):
			ch <- "Enter room name:"
			isCreating = true
			continue
		}

		cast.Messages <- fmt.Sprintf("%s: %s", name, msg)
	}

	cast.Leaving <- ch
	cast.Messages <- fmt.Sprintf("%s exit the chat!\n", name)
	conn.Close()
}
