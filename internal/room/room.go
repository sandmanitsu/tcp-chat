package room

import (
	"fmt"
	"strings"
	"sync"
	"tcp-chat/internal/client"
)

const (
	MainRoomName string = "main"
)

type Room struct {
	Name     string
	Clients  map[client.ClientConn]*client.Client
	Mu       *sync.RWMutex
	Commands []string
}

func GetUsersInRoom(r Room) string {
	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("\tusers in room %s:\n", r.Name))
	for _, client := range r.Clients {
		text.WriteString(fmt.Sprintf("\t-> %s\n", client.Name))
	}

	return text.String()
}

func GetRoomsNames(rooms map[string]Room) string {
	text := strings.Builder{}
	for _, room := range rooms {
		text.WriteString(fmt.Sprintf("\t-> %s\n", room.Name))
	}

	return text.String()
}
