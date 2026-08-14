package orders

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createOrderRequest struct {
	CustomerName string     `json:"customer_name"`
	DueDate      string     `json:"due_date"`
	LineItems    []LineItem `json:"line_items"`
}

type controller struct {
	service *Service
}

func NewController(service *Service) *controller {
	return &controller{service: service}
}

func (h *controller) CreateOrder(c *gin.Context) {
	var req createOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid due_date",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err = h.service.CreateOrder(
		c.Request.Context(),
		userID,
		req.CustomerName,
		dueDate,
		req.LineItems,
	)

	if errors.Is(err, ErrInvalidOrder) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create order",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "order created",
	})
}

func (h *controller) ListOrders(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	orders, err := h.service.ListOrders(
		c.Request.Context(),
		userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func (h *controller) GetOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	order, err := h.service.GetOrder(
		c.Request.Context(),
		userID,
		orderID,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *controller) UpdateOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid due_date",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err = h.service.UpdateOrder(
		c.Request.Context(),
		userID,
		orderID,
		req.CustomerName,
		dueDate,
		req.LineItems,
	)

	if errors.Is(err, ErrInvalidOrder) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order",
		})
		return
	}

	if errors.Is(err, ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	if errors.Is(err, ErrOrderHasPayments) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "order cannot be modified after payment",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update order",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "order updated",
	})
}

func (h *controller) DeleteOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	err = h.service.DeleteOrder(
		c.Request.Context(),
		userID,
		orderID,
	)

	if errors.Is(err, ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	if errors.Is(err, ErrOrderHasPayments) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "order cannot be deleted after payment",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete order",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
