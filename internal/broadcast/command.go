package broadcast

import (
	"fmt"
	"strings"
)

const (
	// Command name
	GetUsers   = "--users"
	CreateRoom = "--create"

	Help = "--help"

	// Command description
	GetUsersDesc    = "get all users into the room"
	CreateRoomsDesc = "create new room"
)

func GetCommandsDescription(commands []string) string {
	desc := map[string]string{
		GetUsers:   GetUsersDesc,
		CreateRoom: CreateRoomsDesc,
	}

	text := strings.Builder{}
	for _, command := range commands {
		text.WriteString(fmt.Sprintf("\t%s  %s\n", command, desc[command]))
	}

	return text.String()
}
