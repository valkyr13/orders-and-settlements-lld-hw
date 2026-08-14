package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitializeRoutes(router *gin.Engine, pool *pgxpool.Pool) {

	healthHandler := func(c *gin.Context) {
		if err := pool.Ping(c.Request.Context()); err != nil {
			c.JSON(500, gin.H{
				"status": "DB unhealthy",
				"error":  err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"status": "ok",
		})
	}

	router.GET("/health", healthHandler)
}
