package main

import (
	"context"
	"log"
	"net/http"
	"orders-and-settlements-lld-hw/backend/config"
	"orders-and-settlements-lld-hw/backend/repository"
	"orders-and-settlements-lld-hw/backend/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	cfg := config.Load()

	pool, err := repository.NewPool(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	router := gin.Default()

	routes.InitializeRoutes(router, pool)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}



}