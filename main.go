package main

import (
	"backend/bootstrap"
	"backend/database"
	"backend/routes"
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	databaseURL := os.Getenv("DB_URL")
	db, err := database.ConnectDB(databaseURL)

	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())

	if err != nil {
		panic(err)
	}

	app := bootstrap.NewApplication(db)
	
	r := gin.Default()

	routes.SetupRoutes(r,app)
	
	r.Run()
}
