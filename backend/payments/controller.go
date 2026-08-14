package payments

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type controller struct {
	service *Service
}

func NewController(service *Service) *controller {
	return &controller{service: service}
}

type createPaymentRequest struct {
	AmountCents int64   `json:"amount_cents"`
	PaymentDate string  `json:"payment_date"`
	Note        *string `json:"note"`
}

func (h *controller) CreatePayment(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	var req createPaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	paymentDate, err := time.Parse(time.RFC3339, req.PaymentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid payment_date",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	payment, err := h.service.CreatePayment(
		c.Request.Context(),
		userID,
		orderID,
		req.AmountCents,
		paymentDate,
		req.Note,
	)

	if errors.Is(err, ErrInvalidPayment) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "payment amount must be greater than zero",
		})
		return
	}

	if errors.Is(err, ErrOrderNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "order not found",
		})
		return
	}

	if errors.Is(err, ErrPaymentExceedsDue) {
		c.JSON(http.StatusConflict, gin.H{
			"error": "payment exceeds amount due",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create payment",
		})
		return
	}

	c.JSON(http.StatusCreated, payment)
}

func (h *controller) ListPayments(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid order id",
		})
		return
	}

	userID := c.MustGet("user_id").(uuid.UUID)

	payments, err := h.service.ListPayments(
		c.Request.Context(),
		userID,
		orderID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch payments",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": payments,
	})
}
