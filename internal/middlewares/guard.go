package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"github.com/rickferrdev/gongo-simple-auth/internal/platform"
	"github.com/rickferrdev/gongo-simple-auth/internal/utils"
)

type GuardMiddleware struct {
	tokenizer platform.Tokenizer
}

func NewGuardMiddleware(router *gin.RouterGroup, tokenizer platform.Tokenizer) GuardMiddleware {
	guard := GuardMiddleware{
		tokenizer: tokenizer,
	}

	router.Use(guard.Guard())

	return guard
}

func (u *GuardMiddleware) Guard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		get := ctx.GetHeader("Authorization")
		parts := strings.Split(get, " ")

		if parts[0] != "Bearer" || parts[1] == "" {
			utils.CaptureHTTP(ctx, domain.ErrUnauthorized)
			return
		}

		claims, err := u.tokenizer.ValidateToken(parts[1])
		if err != nil {
			utils.CaptureHTTP(ctx, domain.ErrUnauthorized)
			return
		}

		if len(claims.Subject) == 0 {
			utils.CaptureHTTP(ctx, domain.ErrUnauthorized)
		}

		ctx.Set("user_id", claims.Subject)
		ctx.Next()
	}
}
