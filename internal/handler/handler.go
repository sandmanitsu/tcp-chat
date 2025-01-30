package handler

import (
	"bufio"
	"fmt"
	"net"
	"tcp-chat/internal/broadcast"
	"tcp-chat/internal/client"
)

func HandleConn(conn net.Conn, cast broadcast.Broadcast) {
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

	ch <- fmt.Sprintf("Welcome, %s!\n", name)
	cast.Entering <- broadcast.EnteringData{
		Chan: ch,
		Client: client.Client{
			Name: name,
		},
	}
	cast.Messages <- fmt.Sprintf("%s enter the chat!\n", name)

	for input.Scan() {
		cast.Messages <- fmt.Sprintf("%s: %s", name, input.Text())
	}

	cast.Leaving <- ch
	cast.Messages <- fmt.Sprintf("%s exit the chat!\n", name)
	conn.Close()
}

func writer(conn net.Conn, ch chan string) {
	for msg := range ch {
		fmt.Fprintf(conn, "%s\n", msg)
	}
}
