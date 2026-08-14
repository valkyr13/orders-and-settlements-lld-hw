package main

import (
	"context"
	"log"
	"net/http"
	"orders-and-settlements-lld-hw/backend/config"
	"orders-and-settlements-lld-hw/backend/repository"
	"orders-and-settlements-lld-hw/backend/routes"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	pool, err := repository.NewPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.InitializeRoutes(router, pool)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
