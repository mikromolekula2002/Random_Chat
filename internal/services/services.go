package services

import (
	roommanager "mw-chat-websocket/internal/roomManager"
	wseventmanager "mw-chat-websocket/internal/wsEventManager"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ISocketService interface {
	ConnectSocket(ctx *gin.Context) error
	WebSocketManager(ctx *gin.Context, conn *websocket.Conn, userID string) error
}

type Service struct {
	ISocketService
}

func NewService(roomManager *wseventmanager.EventManager, eventManager *roommanager.RoomManager) *Service {
	return &Service{
		ISocketService: NewSocketService(roomManager, eventManager),
	}
}
