package routes

import "github.com/gin-gonic/gin"



func InitializeRoutes(router *gin.Engine) {
	healthHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}
	
	router.GET("/health", healthHandler)
}
