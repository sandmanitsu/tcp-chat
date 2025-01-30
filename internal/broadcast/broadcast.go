package broadcast

import (
	"tcp-chat/internal/client"
	"tcp-chat/internal/room"
)

type Broadcast struct {
	Entering chan EnteringData
	Messages chan string
	Leaving  chan client.ClientConn
	Room     room.Room
}

type EnteringData struct {
	Chan   client.ClientConn
	Client client.Client
}

func NewBroadcast(room room.Room) Broadcast {
	return Broadcast{
		Entering: make(chan EnteringData),
		Messages: make(chan string),
		Leaving:  make(chan client.ClientConn),
		Room:     room,
	}
}

func (b *Broadcast) Broadcast() {
	for {
		select {
		case msg := <-b.Messages:
			for cliConn := range b.Room.Clients {
				cliConn <- msg
			}
		case cli := <-b.Entering:
			b.Room.Mu.Lock()
			b.Room.Clients[cli.Chan] = cli.Client
			b.Room.Mu.Unlock()
		case cli := <-b.Leaving:
			b.Room.Mu.Lock()
			delete(b.Room.Clients, cli)
			b.Room.Mu.Unlock()
		}
	}
}
