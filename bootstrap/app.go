package bootstrap

import (
	"backend/handlers"
	"backend/repositories"
	"backend/services"
	"backend/ws"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	MessageHandler *handlers.MessageHandler
}

func NewApplication(db *pgxpool.Pool) *Application {

	manager := ws.NewConnectionManager()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	messageService := services.NewMessageService(messageRepo)
	
	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	messageHandler := handlers.NewMessageHandler(messageService,manager)

	

	return &Application{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		MessageHandler: messageHandler,
	}
}