package bootstrap

import (
	"backend/handlers"
	"backend/repositories"
	"backend/services"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	AuthHandler *handlers.AuthHandler
	UserHandler *handlers.UserHandler
}

func NewApplication(db *pgxpool.Pool) *Application {

	// Repositories
	userRepo := repositories.NewUserRepository(db)

	// Services
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	return &Application{
		AuthHandler: authHandler,
		UserHandler: userHandler,
	}
}