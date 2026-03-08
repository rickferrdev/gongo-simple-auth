package auth

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"github.com/rickferrdev/gongo-simple-auth/internal/utils"
)

type AuthHandler struct {
	service AuthService
}

type RequestLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RequestRegisterDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ResponseLoginDTO struct {
	ID    string `json:"id"`
	Token string `json:"token"`
}

type ResponseRegisterDTO struct {
	ID string `json:"id"`
}

func NewAuthHandler(router *gin.RouterGroup, service AuthService) AuthHandler {
	auth := AuthHandler{
		service: service,
	}

	router.POST("/login", auth.Login)
	router.POST("/register", auth.Register)

	return auth
}

// ErrUnauthorized       = errors.New("unauthorized")
// ErrUserNotFound       = errors.New("user not found")
// ErrUserAlreadyExists  = errors.New("user already exists")
// ErrInvalidCredentials = errors.New("invalid credentials")

// ErrTokenExpired   = errors.New("token has expired")
// ErrTokenInvalid   = errors.New("token is invalid")
// ErrTokenMalformed = errors.New("token format is incorrect")

// ErrInternal = errors.New("internal server error")
// ErrTimeout  = errors.New("operation timed out")
// ErrDatabase = errors.New("database provider error")

func (u *AuthHandler) capture(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrUnauthorized):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": ErrUserNotFound.Error()})
	case errors.Is(err, domain.ErrUserAlreadyExists):
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": domain.ErrUserAlreadyExists.Error()})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrInvalidCredentials.Error()})
	case errors.Is(err, domain.ErrTimeout):
		c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{"error": domain.ErrTimeout.Error()})
	case errors.Is(err, domain.ErrTokenMalformed):
	case errors.Is(err, domain.ErrTokenInvalid):
	case errors.Is(err, domain.ErrTokenExpired):
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": domain.ErrTokenInvalid.Error()})
	default:
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": domain.ErrInternal.Error()})
	}
}

func (u *AuthHandler) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var body RequestRegisterDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.CaptureHTTP(c, domain.ErrBadRequest)
		return
	}

	input := RegisterInput{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
	}

	output, err := u.service.Register(ctx, input)
	if err != nil {
		utils.CaptureHTTP(c, err)
		return
	}

	c.JSON(http.StatusCreated, ResponseRegisterDTO{ID: output.ID})
}

func (u *AuthHandler) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	var body RequestLoginDTO
	if err := c.ShouldBindJSON(&body); err != nil {
		utils.CaptureHTTP(c, domain.ErrBadRequest)
		return
	}

	input := LoginInput{
		Email:    body.Email,
		Password: body.Password,
	}

	output, err := u.service.Login(ctx, input)
	if err != nil {
		utils.CaptureHTTP(c, err)
		return
	}

	c.JSON(http.StatusOK, ResponseLoginDTO{ID: output.ID, Token: output.Token})
}
