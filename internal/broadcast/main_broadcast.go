package broadcast

import (
	"fmt"
	"strings"
	"sync"
	"tcp-chat/internal/client"
	"tcp-chat/internal/room"
)

type MainBroadcast struct {
	Entering chan EnteringData
	Messages chan string
	Leaving  chan client.ClientConn
	Creating chan CreatingData
	Room     *room.MainRoom
	Commands []string
	Mu       *sync.Mutex
}

type CreatingData struct {
	Chan     client.ClientConn
	Client   client.Client
	RoomName string
}

func NewMainBroardcast(room *room.MainRoom) *MainBroadcast {
	return &MainBroadcast{
		Entering: make(chan EnteringData),
		Messages: make(chan string),
		Leaving:  make(chan client.ClientConn),
		Creating: make(chan CreatingData),
		Room:     room,
		Commands: []string{GetUsers, CreateRoom},
		Mu:       new(sync.Mutex),
	}
}

func (m *MainBroadcast) Broadcast() {
	for {
		select {
		case msg := <-m.Messages:
			for cliConn := range m.Room.Clients {
				cliConn <- msg
			}
		case cli := <-m.Entering:
			m.Room.Mu.Lock()
			m.Room.Clients[cli.Chan] = cli.Client
			m.Room.Mu.Unlock()
		case cli := <-m.Leaving:
			m.Room.Mu.Lock()
			delete(m.Room.Clients, cli)
			m.Room.Mu.Unlock()
		case r := <-m.Creating:
			_, ok := m.Room.Rooms[r.RoomName]
			if ok {
				r.Chan <- fmt.Sprintf("room with name: %s already exist", r.RoomName)

				continue
			}

			m.Room.Mu.Lock()

			newRoom := room.NewRoom()
			newRoom.Clients[r.Chan] = r.Client
			broadcast := NewBroadcast(newRoom)
			go broadcast.Broadcast()

			m.Room.Rooms[r.RoomName] = newRoom
			delete(m.Room.Clients, r.Chan)

			m.Room.Mu.Unlock()

			r.Chan <- fmt.Sprintf("\twelcome to romm %s", r.RoomName)
		}
	}
}

func (m *MainBroadcast) GetUsers() string {
	text := strings.Builder{}
	text.WriteString("\tusers:\n")
	for _, client := range m.Room.Clients {
		text.WriteString(fmt.Sprintf("\t %s\n", client.Name))
	}

	return text.String()
}
