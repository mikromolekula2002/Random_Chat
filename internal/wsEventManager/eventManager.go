package wseventmanager

import (
	"fmt"
	"mw-chat-websocket/internal/schemas"

	"github.com/gorilla/websocket"
)

type EventManager struct {
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (em *EventManager) SendMessageToRoom(messageData *schemas.MessageRequest, room *schemas.Room) error {
	op := "SendMessageToRoom"

	if len(room.Connections) == 0 {
		return fmt.Errorf("%s: There are no users in the room - %s", op, messageData.Room)
	}

	message := &schemas.MessageResponse{
		Action:        "sendMessage",
		SendlerUserID: messageData.UserID,
		Message:       messageData.Message,
	}

	// Sending an event to both users
	for _, clientConn := range room.Connections {
		if err := clientConn.WriteJSON(message); err != nil {
			return fmt.Errorf("%s: %v", op, err)
		}
	}

	return nil
}

func (em *EventManager) SendRoomUpdate(conn *websocket.Conn, roomID, userID string) error {
	op := "SendRoomUpdate"

	message := &schemas.RoomUpdateResponse{
		Action: "roomUpdate",
		RoomID: roomID,
		UserID: userID,
	}

	err := conn.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}

func (em *EventManager) SendSearchResult(conn *websocket.Conn, result string) error {
	op := "SendSearchResult"

	message := &schemas.SearchResultResponse{
		Action: "searchResult",
		Result: result,
	}

	err := conn.WriteJSON(message)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	return nil
}
