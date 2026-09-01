package routes

import (
	"backend/bootstrap"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(r *gin.Engine, app *bootstrap.Application){

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	r.POST("/register", app.AuthHandler.Register)
	r.POST("/login", app.AuthHandler.Login)
	
	protected.GET("/users", app.UserHandler.GetUsers)

	protected.POST("/message",app.MessageHandler.SendMessage)

	protected.GET("/message/:otherUserID",app.MessageHandler.GetConversation)

	protected.GET("/ws", app.MessageHandler.WebSocket)

}