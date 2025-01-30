package room

import (
	"sync"
	"tcp-chat/internal/client"
)

type Room struct {
	Clients map[client.ClientConn]client.Client
	Mu      *sync.RWMutex
}

func NewRoom() Room {
	return Room{
		Clients: make(map[client.ClientConn]client.Client),
		Mu:      new(sync.RWMutex),
	}
}
