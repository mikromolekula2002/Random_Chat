package roommanager

import (
	"fmt"
	"mw-chat-websocket/internal/schemas"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	connectRoom = "Another user has been found, let's chat with him"
	createRoom  = "Waiting for other user"
	leftRoom    = "The user has left the chat, we are waiting for a new user"
)

var Rooms = make(map[string]*schemas.Room)

type RoomManager struct {
}

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}

func (rm *RoomManager) FindOrCreateRoom(userID string, conn *websocket.Conn, avoidRoomID string) (*schemas.Room, string, error) {
	op := "FindOrCreateRoom"

	// We are looking for an available room, excluding the room user just left.
	for _, room := range Rooms {
		if !room.IsFull && room.ID != avoidRoomID {
			room, searchResult, err := rm.ConnectToRoom(conn, userID, room)
			if err != nil {
				return nil, "", fmt.Errorf("%s: %v", op, err)
			}

			return room, searchResult, nil
		}
	}

	// If there are no rooms available, we check if the old one still exists
	if avoidRoomID != "" {
		if room, exists := Rooms[avoidRoomID]; exists && !room.IsFull {
			room, searchResult, err := rm.ConnectToRoom(conn, userID, room)
			if err != nil {
				return nil, "", fmt.Errorf("%s: %v", op, err)
			}

			return room, searchResult, nil
		}
	}

	// Create a new room if there are no rooms available and the previous one is also unavailable
	room, searchResult, err := rm.CreateNewRoom(conn, userID)
	if err != nil {
		return nil, "", fmt.Errorf("%s: %v", op, err)
	}

	return room, searchResult, nil
}

func (rm *RoomManager) ConnectToRoom(conn *websocket.Conn, userID string, room *schemas.Room) (*schemas.Room, string, error) {

	room.Connections[userID] = conn
	if len(room.Connections) >= 2 {
		room.IsFull = true
	}

	return room, connectRoom, nil
}

func (rm *RoomManager) CreateNewRoom(conn *websocket.Conn, userID string) (*schemas.Room, string, error) {

	newRoomID := uuid.New().String()
	newRoom := &schemas.Room{
		ID:          newRoomID,
		Connections: make(map[string]*websocket.Conn),
		IsFull:      false,
	}
	newRoom.Connections[userID] = conn
	Rooms[newRoomID] = newRoom

	return newRoom, createRoom, nil
}

func (rm *RoomManager) RemoveConnectionFromRoom(room *schemas.Room, userID string) (string, error) {

	delete(Rooms[room.ID].Connections, userID)

	if len(room.Connections) == 0 {
		delete(Rooms, room.ID)
	} else {
		room.IsFull = false
	}

	return leftRoom, nil
}

func (rm *RoomManager) FindRoomByUserID(userID string) (*schemas.Room, bool) {
	for _, room := range Rooms {
		if _, exists := room.Connections[userID]; exists {
			return room, true
		}
	}
	return nil, false
}

func (rm *RoomManager) RoomExists(roomID string) (*schemas.Room, error) {
	op := "RoomExists"

	room, exists := Rooms[roomID]
	if !exists {
		return nil, fmt.Errorf("%s: This room doesn't exist", op)
	}

	return room, nil
}
