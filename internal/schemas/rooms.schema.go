package schemas

import "github.com/gorilla/websocket"

// the action field accepts the following. parameters: sendMessage, researchRoom
type MessageRequest struct {
	Action  string `json:"action"`
	UserID  string `json:"userID"`
	Room    string `json:"room"`
	Message string `json:"message"`
}

type MessageResponse struct {
	Action        string `json:"action"`
	SendlerUserID string `json:"userID"`
	Message       string `json:"message"`
}

type RoomUpdateResponse struct {
	Action string `json:"action"`
	RoomID string `json:"roomID"`
	UserID string `json:"userID"`
}

type SearchResultResponse struct {
	Action string `json:"action"`
	Result string `json:"result"`
}

// Stores a users websocket conn by type: [userID]ws.conn
type Room struct {
	ID          string
	Connections map[string]*websocket.Conn
	IsFull      bool
}
