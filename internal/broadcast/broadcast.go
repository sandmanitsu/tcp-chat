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
			room.MainRoomName: {
				Name:     room.MainRoomName,
				Clients:  make(map[client.ClientConn]*client.Client),
				Mu:       new(sync.RWMutex),
				Commands: []string{room.GetUsers, room.CreateRoom, room.EnterChat},
			},
		},
		Mu: new(sync.Mutex),
	}
}

func (b Broadcast) Broadcast() {
	for {
		select {
		case msg := <-b.Messages:
			_, ok := b.Room[msg.RoomName]
			if !ok {
				log.Printf("error send message: %s", msg.Msg)

				continue // todo. Если не нашлась комната, отправить юзеру ошибку в чат
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

			for cliConn := range b.Room[e.RoomName].Clients {
				cliConn <- fmt.Sprintf("user %s enter to chat %s", e.Client.Name, e.RoomName)
			}
		case l := <-b.Leaving:
			_, ok := b.Room[l.Client.CurrRoom]
			if !ok {
				l.Chan <- fmt.Sprintf("error leaving chat %s", l.Client.CurrRoom)

				continue
			}

			oldChat := l.Client.CurrRoom

			b.Mu.Lock()
			delete(b.Room[l.Client.CurrRoom].Clients, l.Chan)

			if !l.IsExit {
				b.Room[room.MainRoomName].Clients[l.Chan] = l.Client
				l.Client.CurrRoom = room.MainRoomName

				for cliConn := range b.Room[room.MainRoomName].Clients {
					cliConn <- fmt.Sprintf("user %s enter to chat %s", l.Client.Name, room.MainRoomName)
				}
			}
			b.Mu.Unlock()

			for cliConn := range b.Room[oldChat].Clients {
				cliConn <- fmt.Sprintf("user %s leave this chan", l.Client.Name)
			}
		case r := <-b.Creating:
			b.Room[r.RoomName] = room.Room{
				Name:     r.RoomName,
				Clients:  make(map[client.ClientConn]*client.Client),
				Mu:       new(sync.RWMutex),
				Commands: []string{room.GetUsers, room.LeaveChat},
			}

			b.Mu.Lock()
			delete(b.Room[r.Client.CurrRoom].Clients, r.Chan)
			b.Mu.Unlock()

			_, ok := b.Room[r.RoomName]
			if !ok {
				r.Chan <- fmt.Sprintf("can't connect to room %s", r.RoomName)

				continue
			}

			b.Mu.Lock()
			b.Room[r.RoomName].Clients[r.Chan] = r.Client
			b.Mu.Unlock()

			r.Client.CurrRoom = r.RoomName

			for cliConn := range b.Room[r.RoomName].Clients {
				cliConn <- fmt.Sprintf("user %s enter to chat %s", r.Client.Name, r.RoomName)
			}
		}
	}
}
