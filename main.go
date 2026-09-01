package main

import (
	"backend/bootstrap"
	"backend/database"
	"backend/routes"
	"context"
	"fmt"
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

	redisClient, err := database.ConnectRedis()

	if err != nil {
		panic(err)
	}
	defer redisClient.Close()

	fmt.Println("Redis connected successfully")

	app := bootstrap.NewApplication(db,redisClient)

	r := gin.Default()

	routes.SetupRoutes(r, app)

	r.Run()
}
