package chat

import (
	"errors"
	"sync"
)

type ChatServer struct {
	Address  string
	Rooms    []*ChatRoom
	roomsMut sync.RWMutex
}

type ChatRoom struct {
	Name string
}

type ChatClient struct {
	Address string
}

var (
	ErrRoomAlreadyExists = errors.New("room already exists")
)

func NewChatServer(address string) *ChatServer {
	return &ChatServer{
		Address: address,
		Rooms:   []*ChatRoom{},
	}
}

func (cs *ChatServer) GetRoom(name string) *ChatRoom {
	cs.roomsMut.RLock()
	defer cs.roomsMut.RUnlock()
	for _, room := range cs.Rooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}

func (cs *ChatServer) ListRooms() []*ChatRoom {
	cs.roomsMut.RLock()
	defer cs.roomsMut.RUnlock()
	return cs.Rooms
}

func (cs *ChatServer) AddRoom(name string) (*ChatRoom, error) {
	room := cs.GetRoom(name)
	if room != nil {
		return nil, ErrRoomAlreadyExists
	}
	cs.roomsMut.Lock()
	defer cs.roomsMut.Unlock()
	newRoom := &ChatRoom{Name: name}
	cs.Rooms = append(cs.Rooms, newRoom)
	return newRoom, nil
}

// type ChatServer struct {
// 	Name  string
// 	rooms []ChatRoom
// }

// type ChatRoom struct {
// 	ID int
// }

// var chatServers = map[string]*ChatServer{}
// var chatServerMut = &sync.RWMutex{}

// func newChatServer(name string) *ChatServer {
// 	return &ChatServer{
// 		Name:  name,
// 		rooms: []ChatRoom{},
// 	}
// }

// func getChatServer(serverAddr string) *ChatServer {
// 	chatServerMut.RLock()
// 	defer chatServerMut.RUnlock()

// 	if chatServer, ok := chatServers[serverAddr]; ok {
// 		return chatServer
// 	}

// 	chatServerMut.RUnlock()
// 	chatServerMut.Lock()
// 	defer chatServerMut.Unlock()

// 	// Double-check in case another goroutine created it
// 	if chatServer, ok := chatServers[serverAddr]; ok {
// 		return chatServer
// 	}

// 	chatServer := NewChatServer()
// 	chatServers[serverAddr] = chatServer
// 	return chatServer
// }

// func getRoom() {

// }
