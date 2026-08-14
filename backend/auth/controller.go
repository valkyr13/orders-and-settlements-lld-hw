package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type controller struct {
	service *Service
}

func NewController(service *Service) *controller {
	return &controller{service: service}
}

func (c *controller) Signup(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.service.Signup(ctx.Request.Context(), req.Email, req.Password); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *controller) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := c.service.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := c.service.CreateSession(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie(
		"session_id",
		sessionID.String(),
		86400,
		"/",
		"",
		false,
		true,
	)

	ctx.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
	})
}


func (h *controller) Me(c *gin.Context) {
	userID := c.MustGet("user_id").(uuid.UUID)

	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
	})
}