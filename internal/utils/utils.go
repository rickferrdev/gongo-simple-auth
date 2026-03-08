package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
)

func CaptureHTTP(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrBadRequest):
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": domain.ErrBadRequest.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": domain.ErrUnauthorized.Error()})
	case errors.Is(err, domain.ErrUserNotFound):
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": domain.ErrUserNotFound.Error()})
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
