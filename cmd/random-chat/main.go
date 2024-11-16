package main

import (
	"log"
	"mw-chat-websocket/internal/config"
	"mw-chat-websocket/internal/controllers"
	roommanager "mw-chat-websocket/internal/roomManager"
	"mw-chat-websocket/internal/server"
	wseventmanager "mw-chat-websocket/internal/wsEventManager"
)

// @title Random-Chat API
// @version 1.0
// @description API для анонимного чата с использованием WebSocket.
// @BasePath /random
func main() {
	newConfig, err := config.LoadConfig("")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	eventManager := wseventmanager.NewEventManager()
	roomManager := roommanager.NewRoomManager(eventManager)

	controllers := controllers.NewController(eventManager, roomManager)

	newServer := server.NewServer(&newConfig)
	newServer.SetRoutes(controllers)

	log.Fatal(newServer.GinServer.Run(":" + newConfig.ServerPort))
}
