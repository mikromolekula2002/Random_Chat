package controllers

import (
	roommanager "mw-chat-websocket/internal/roomManager"
	"mw-chat-websocket/internal/services"
	wseventmanager "mw-chat-websocket/internal/wsEventManager"
)

type Controller struct {
	SocketController *SocketController
}

func NewController(roomManager *wseventmanager.EventManager, eventManager *roommanager.RoomManager) *Controller {
	return &Controller{
		SocketController: NewSocketController(*services.NewSocketService(roomManager, eventManager)),
	}
}
