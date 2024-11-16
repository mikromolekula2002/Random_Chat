package roommanager

import (
	"fmt"
	"mw-chat-websocket/internal/schemas"
	wseventmanager "mw-chat-websocket/internal/wsEventManager"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	connectRoom = "Another user has been found, let's chat with him"
	createRoom  = "Waiting for other user"
	leftRoom    = "The user has left the chat, we are waiting for a new user"
)

var Rooms = make(map[string]*schemas.Room)

type IEventManager interface {
	SendMessageToRoom(messageData *schemas.MessageRequest, room *schemas.Room) error
	SendRoomUpdate(conn *websocket.Conn, roomID, userID string) error
	SendSearchResult(conn *websocket.Conn, result string) error
}

type RoomManager struct {
	eventManager IEventManager
}

func NewRoomManager(eventManager *wseventmanager.EventManager) *RoomManager {
	return &RoomManager{eventManager: eventManager}
}

func (rm *RoomManager) FindOrCreateRoom(userID string, conn *websocket.Conn, avoidRoomID string) (*schemas.Room, error) {
	op := "FindOrCreateRoom"

	// We are looking for an available room, excluding the room user just left.
	for _, room := range Rooms {
		if !room.IsFull && room.ID != avoidRoomID {
			room, err := rm.ConnectToRoom(conn, userID, room)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", op, err)
			}

			return room, nil
		}
	}

	// If there are no rooms available, we check if the old one still exists
	if avoidRoomID != "" {
		if room, exists := Rooms[avoidRoomID]; exists && !room.IsFull {
			room, err := rm.ConnectToRoom(conn, userID, room)
			if err != nil {
				return nil, fmt.Errorf("%s: %v", op, err)
			}

			return room, nil
		}
	}

	// Create a new room if there are no rooms available and the previous one is also unavailable
	room, err := rm.CreateNewRoom(conn, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return room, nil
}

func (rm *RoomManager) ConnectToRoom(conn *websocket.Conn, userID string, room *schemas.Room) (*schemas.Room, error) {
	op := "ConnectToRoom"

	room.Connections[userID] = conn
	if len(room.Connections) >= 2 {
		room.IsFull = true
	}

	for _, c := range room.Connections {
		err := rm.eventManager.SendSearchResult(c, connectRoom)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", op, err)
		}
	}

	err := rm.eventManager.SendRoomUpdate(conn, room.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return room, nil
}

func (rm *RoomManager) CreateNewRoom(conn *websocket.Conn, userID string) (*schemas.Room, error) {
	op := "CreateNewRoom"

	newRoomID := uuid.New().String()
	newRoom := &schemas.Room{
		ID:          newRoomID,
		Connections: make(map[string]*websocket.Conn),
		IsFull:      false,
	}
	newRoom.Connections[userID] = conn
	Rooms[newRoomID] = newRoom

	err := rm.eventManager.SendSearchResult(conn, createRoom)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	err = rm.eventManager.SendRoomUpdate(conn, newRoom.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", op, err)
	}

	return newRoom, nil
}

func (rm *RoomManager) RemoveConnectionFromRoom(room *schemas.Room, userID string) error {
	op := "RemoveConnectionFromRoom"

	delete(Rooms[room.ID].Connections, userID)

	for _, wsConn := range Rooms[room.ID].Connections {
		if err := rm.eventManager.SendSearchResult(wsConn, leftRoom); err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	if len(room.Connections) == 0 {
		delete(Rooms, room.ID)
	} else {
		room.IsFull = false
	}

	return nil
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
