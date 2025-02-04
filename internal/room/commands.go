package room

import (
	"fmt"
	"strings"
)

const (
	// Command name
	GetUsers   = "--users"
	CreateRoom = "--create"
	EnterChat  = "--enter"
	LeaveChat  = "--leave"

	Help = "--help"

	// Command description
	GetUsersDesc    = "get all users into the room"
	CreateRoomsDesc = "create new room"
	EnterChatDesc   = "enter to existing chat"
	LeaveChatDesc   = "leave current chat if current isn't 'main'"
)

func GetCommandsDescription(commands []string) string {
	desc := map[string]string{
		GetUsers:   GetUsersDesc,
		CreateRoom: CreateRoomsDesc,
		EnterChat:  EnterChatDesc,
		LeaveChat:  LeaveChatDesc,
	}

	text := strings.Builder{}
	for _, command := range commands {
		text.WriteString(fmt.Sprintf("\t%s  %s\n", command, desc[command]))
	}

	return text.String()
}
