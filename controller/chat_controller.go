package controller

import "github.com/gin-gonic/gin"

type ChatController interface {
	GetMessageHistory(ctx *gin.Context)
	WSHandler(ctx *gin.Context)
	GetUser(ctx *gin.Context)
}
