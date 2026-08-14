package routes

import (
	"orders-and-settlements-lld-hw/backend/auth"
	"orders-and-settlements-lld-hw/backend/middleware"
	"orders-and-settlements-lld-hw/backend/orders"
	"orders-and-settlements-lld-hw/backend/payments"

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

	orderRepo := orders.NewRepository(pool)
	orderService := orders.NewService(orderRepo)
	orderController := orders.NewController(orderService)

	paymentRepo := payments.NewRepository(pool)
	paymentService := payments.NewService(paymentRepo)
	paymentController := payments.NewController(paymentService)

	authRepository := auth.NewRepository(pool)
	authService := auth.NewService(authRepository)
	authController := auth.NewController(authService)


	router.GET("/health", healthHandler)

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/signup", authController.Signup)
		authRoutes.POST("/login", authController.Login)
	}

	protected := router.Group("/")
	protected.Use(middleware.RequireAuth(authRepository))
	{
		protected.POST("/orders", orderController.CreateOrder)
		protected.GET("/orders", orderController.ListOrders)
		protected.GET("/orders/:id", orderController.GetOrder)
		protected.PUT("/orders/:id", orderController.UpdateOrder)
		protected.DELETE("/orders/:id", orderController.DeleteOrder)

		protected.POST("/orders/:id/payments", paymentController.CreatePayment)
		protected.GET("/orders/:id/payments", paymentController.ListPayments)
		protected.GET("/auth/me", authController.Me)

	}
}
