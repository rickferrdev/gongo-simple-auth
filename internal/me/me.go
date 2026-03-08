package me

import (
	"github.com/gin-gonic/gin"
	"github.com/rickferrdev/gongo-simple-auth/internal/domain"
	"github.com/rickferrdev/gongo-simple-auth/internal/utils"
)

type MeHandler struct{}

func NewMeHandler(router *gin.RouterGroup) MeHandler {
	me := MeHandler{}

	router.GET("/me", me.Me)
	return me
}

func (u *MeHandler) Me(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.CaptureHTTP(c, domain.ErrUnauthorized)
		return
	}

	c.JSON(200, gin.H{"data": userID})
}
