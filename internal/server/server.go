package server

import (
	"mw-chat-websocket/internal/config"
	"mw-chat-websocket/internal/controllers"
	"net/http"

	_ "mw-chat-websocket/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	GinServer *gin.Engine
	Config    *config.Config
}

func NewServer(cfg *config.Config) *Server {
	server := gin.Default()

	server.LoadHTMLGlob("templates/*")

	server.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "The way APi is working fine"})
	})

	if cfg.EnvType != "prod" {
		server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return &Server{
		GinServer: server,
		Config:    cfg,
	}
}

func (server *Server) SetRoutes(controller *controllers.Controller) {
	mwChatWebsocket := server.GinServer.Group("/random")
	{
		mwChatWebsocket.GET("/ws", controller.SocketController.ConnectSocket)
		mwChatWebsocket.GET("/chat", controller.SocketController.ChatPage)
		mwChatWebsocket.GET("/home", controller.SocketController.HomePage)
	}
}
