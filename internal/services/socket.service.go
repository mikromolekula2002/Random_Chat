package services

import (
	"encoding/json"
	"fmt"
	"log"
	roommanager "mw-chat-websocket/internal/roomManager"
	"mw-chat-websocket/internal/schemas"
	wseventmanager "mw-chat-websocket/internal/wsEventManager"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	emptyAvoidRoomID = ""
)

var SessionPool = make(map[string]*websocket.Conn)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type IEventManager interface {
	SendMessageToRoom(messageData *schemas.MessageRequest, room *schemas.Room) error
	SendRoomUpdate(conn *websocket.Conn, roomID, userID string) error
	SendSearchResult(conn *websocket.Conn, result string) error
}

type IRoomManager interface {
	FindOrCreateRoom(userID string, conn *websocket.Conn, avoidRoomID string) (*schemas.Room, string, error)
	ConnectToRoom(conn *websocket.Conn, userID string, room *schemas.Room) (*schemas.Room, string, error)
	CreateNewRoom(conn *websocket.Conn, userID string) (*schemas.Room, string, error)
	RemoveConnectionFromRoom(room *schemas.Room, userID string) (string, error)
	FindRoomByUserID(userID string) (*schemas.Room, bool)
	RoomExists(roomID string) (*schemas.Room, error)
}

type SocketService struct {
	eventManager IEventManager
	roomManager  IRoomManager
}

func NewSocketService(eventManager *wseventmanager.EventManager, roomManager *roommanager.RoomManager) *SocketService {
	return &SocketService{
		eventManager: eventManager,
		roomManager:  roomManager,
	}
}

func (ss *SocketService) ConnectSocket(ctx *gin.Context) error {
	op := "ConnectSocket"

	userID := uuid.New().String()
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return fmt.Errorf("%s - Error during handshake: %v", op, err)
	}

	SessionPool[userID] = conn

	room, searchResult, err := ss.roomManager.FindOrCreateRoom(userID, conn, emptyAvoidRoomID)
	if err != nil {
		return fmt.Errorf("%s: %v", op, err)
	}

	for _, wsConn := range room.Connections {
		if err := ss.eventManager.SendSearchResult(wsConn, searchResult); err != nil {
			return fmt.Errorf("%s : %v", op, err)
		}
	}

	if err := ss.eventManager.SendRoomUpdate(conn, room.ID, userID); err != nil {
		return fmt.Errorf("%s : %v", op, err)
	}

	go func() {
		if err := ss.WebSocketManager(ctx, conn, userID); err != nil {
			log.Println(err)
		}
	}()

	return nil
}

// Reads and processes events coming from the user
func (ss *SocketService) WebSocketManager(ctx *gin.Context, conn *websocket.Conn, userID string) error {
	op := "WsMessageManager"

	defer func() {
		ss.CloseConnections(conn, userID)
	}()

	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("%s - User has left the chat : %v", op, err)
		}

		var messageData schemas.MessageRequest
		if err := json.Unmarshal(rawMessage, &messageData); err != nil {
			return fmt.Errorf("%s : %v", op, err)
		}

		switch messageData.Action {
		case "sendMessage":
			var room *schemas.Room

			room, err := ss.roomManager.RoomExists(messageData.Room)
			if err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

			if err := ss.eventManager.SendMessageToRoom(&messageData, room); err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

		case "researchRoom":
			room, err := ss.roomManager.RoomExists(messageData.Room)
			if err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

			searchResult, err := ss.roomManager.RemoveConnectionFromRoom(room, userID)
			if err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

			for _, wsConn := range room.Connections {
				if err := ss.eventManager.SendSearchResult(wsConn, searchResult); err != nil {
					return fmt.Errorf("%s : %v", op, err)
				}
			}

			room, searchResult, err = ss.roomManager.FindOrCreateRoom(userID, conn, messageData.Room)
			if err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

			for _, wsConn := range room.Connections {
				if err := ss.eventManager.SendSearchResult(wsConn, searchResult); err != nil {
					return fmt.Errorf("%s : %v", op, err)
				}
			}

			if err := ss.eventManager.SendRoomUpdate(conn, room.ID, userID); err != nil {
				return fmt.Errorf("%s : %v", op, err)
			}

		}
	}
}

// Checks and deletes if there is a connection from the session pool, as well as room connection
func (ss *SocketService) CloseConnections(conn *websocket.Conn, userID string) {
	conn.Close()

	room, found := ss.roomManager.FindRoomByUserID(userID)
	if found {
		searchResult, err := ss.roomManager.RemoveConnectionFromRoom(room, userID)
		if err != nil {
			log.Printf("Error while close user connection : %v", err)
		}

		if err := ss.eventManager.SendSearchResult(conn, searchResult); err != nil {
			log.Printf("Error while close user connection : %v", err)
		}
	}

	delete(SessionPool, userID)
}
