package bootstrap

import (
	"backend/handlers"
	"backend/repositories"
	"backend/services"
	"backend/ws"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
	MessageHandler *handlers.MessageHandler
}

func NewApplication(db *pgxpool.Pool,redis *redis.Client) *Application {

	manager := ws.NewConnectionManager()

	// Repositories
	userRepo := repositories.NewUserRepository(db)
	messageRepo := repositories.NewMessageRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo,redis)
	messageService := services.NewMessageService(messageRepo)
	redisService := services.NewRedisService(redis)
	
	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	
	messageHandler := handlers.NewMessageHandler(messageService,manager, redisService,userService)

	

	return &Application{
		AuthHandler: authHandler,
		UserHandler: userHandler,
		MessageHandler: messageHandler,
	}
}