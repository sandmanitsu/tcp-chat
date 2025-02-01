package broadcast

import (
	"fmt"
	"log"
	"sync"
	"tcp-chat/internal/client"
	"tcp-chat/internal/room"
)

type Broadcast struct {
	Entering chan EnteringData
	Messages chan Message
	Leaving  chan LeavingData
	Creating chan CreatingData
	Room     map[string]room.Room
	Mu       *sync.Mutex
}

type CreatingData struct {
	Chan     client.ClientConn
	Client   *client.Client
	RoomName string
}

type EnteringData struct {
	Chan     client.ClientConn
	Client   *client.Client
	RoomName string
}

type LeavingData struct {
	Chan   client.ClientConn
	Client *client.Client
	IsExit bool
}

type Message struct {
	Msg      string
	RoomName string
}

func NewBroadcast() Broadcast {
	return Broadcast{
		Entering: make(chan EnteringData),
		Messages: make(chan Message),
		Leaving:  make(chan LeavingData),
		Creating: make(chan CreatingData),
		Room: map[string]room.Room{
			room.MainRoomName: room.Room{
				Name:     room.MainRoomName,
				Clients:  make(map[client.ClientConn]*client.Client),
				Mu:       new(sync.RWMutex),
				Commands: []string{room.GetUsers, room.CreateRoom},
			},
		},
		Mu: new(sync.Mutex),
	}
}

func (b *Broadcast) Broadcast() {
	for {
		select {
		case msg := <-b.Messages:
			_, ok := b.Room[msg.RoomName]
			if !ok {
				log.Printf("error send message: %s", msg.Msg)

				continue // todo. Если не нашлось комната, отправить юзеру ошибку в чат
			}

			for cliConn := range b.Room[msg.RoomName].Clients {
				cliConn <- msg.Msg
			}
		case e := <-b.Entering:
			_, ok := b.Room[e.RoomName]
			if !ok {
				e.Chan <- fmt.Sprintf("can't connect to room %s", e.RoomName)

				continue
			}

			b.Mu.Lock()
			b.Room[e.RoomName].Clients[e.Chan] = e.Client
			b.Mu.Unlock()

			e.Client.CurrRoom = e.RoomName

			b.Messages <- Message{
				Msg:      fmt.Sprintf("user %s enter to chat %s", e.Client.Name, e.RoomName),
				RoomName: e.RoomName,
			}
		case l := <-b.Leaving:
			_, ok := b.Room[l.Client.CurrRoom]
			if !ok {
				l.Chan <- fmt.Sprintf("error leaving chan %s", l.Client.CurrRoom)

				continue
			}

			b.Mu.Lock()
			delete(b.Room[l.Client.CurrRoom].Clients, l.Chan)
			b.Mu.Unlock()

			b.Messages <- Message{
				Msg:      fmt.Sprintf("user %s exit the chan %s", l.Client.Name, l.Client.CurrRoom),
				RoomName: l.Client.CurrRoom,
			}

			if !l.IsExit {
				b.Room[room.MainRoomName].Clients[l.Chan] = l.Client
			}
		case r := <-b.Creating:
			b.Room[r.RoomName] = room.Room{
				Name:     r.RoomName,
				Clients:  make(map[client.ClientConn]*client.Client),
				Mu:       new(sync.RWMutex),
				Commands: []string{room.GetUsers},
			}

			b.Mu.Lock()
			delete(b.Room[r.Client.CurrRoom].Clients, r.Chan)

			b.Entering <- EnteringData{
				Chan:     r.Chan,
				Client:   r.Client,
				RoomName: r.RoomName,
			}
		}
	}
}
