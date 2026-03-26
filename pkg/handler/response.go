package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, status int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(status, errorResponse{message})
}
