package controllers

import (
	"log"
	"mw-chat-websocket/internal/services"

	"github.com/gin-gonic/gin"
)

type SocketController struct {
	SocketService services.SocketService
}

func NewSocketController(socketService services.SocketService) *SocketController {
	return &SocketController{SocketService: socketService}
}

// HomePage
// @Summary      Home Page
// @Description  Отображает домашнюю страницу приложения.
// @Tags         Pages
// @Accept       html
// @Produce      html
// @Success      200 {string} string "HTML страницы успешно загружен"
// @Router       / [get]
func (cc *SocketController) HomePage(ctx *gin.Context) {
	ctx.HTML(200, "home.html", nil)
}

// ChatPage
// @Summary      Chat Page
// @Description  Отображает страницу чата приложения.
// @Tags         Pages
// @Accept       html
// @Produce      html
// @Success      200 {string} string "HTML страницы успешно загружен"
// @Router       /chat [get]
func (cc *SocketController) ChatPage(ctx *gin.Context) {
	ctx.HTML(200, "chat.html", nil)
}

// ConnectSocket
// @Summary      Подключение к WebSocket
// @Description  Создает WebSocket соединение для пользователя и обрабатывает события через WebSocket.
// @Tags         WebSocket
// @Produce      json
// @Success      101 {string} string "WebSocket соединение успешно установлено"
// @Failure      400 {object} map[string]string "Ошибка при установке WebSocket соединения"
// @Router       /ws [get]
func (cc *SocketController) ConnectSocket(ctx *gin.Context) {
	err := cc.SocketService.ConnectSocket(ctx)
	if err != nil {
		log.Println(err)
	}
}
