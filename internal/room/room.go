package room

import (
	"sync"
	"tcp-chat/internal/client"
)

type MainRoom struct {
	Rooms   map[string]Room
	Clients map[client.ClientConn]client.Client
	Mu      *sync.RWMutex
}

type Room struct {
	Clients map[client.ClientConn]client.Client
	Mu      *sync.RWMutex
}

func NewMainRooms() *MainRoom {
	return &MainRoom{
		Rooms:   make(map[string]Room),
		Clients: make(map[client.ClientConn]client.Client),
		Mu:      new(sync.RWMutex),
	}
}

func NewRoom() Room {
	return Room{
		Clients: make(map[client.ClientConn]client.Client),
		Mu:      new(sync.RWMutex),
	}
}
