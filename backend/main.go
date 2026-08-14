package main

import (
	"log"
	"net/http"
	"orders-and-settlements-lld-hw/backend/routes"

	"github.com/gin-gonic/gin"
)

func main(){
	//TODO: initialise db connection
	router := gin.Default()

	routes.InitializeRoutes(router)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}



}