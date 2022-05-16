package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"invalid parameter"`
}

func NewHttpError(c *gin.Context, code int, msg string) {
	c.JSON(code, HttpError{Code: code, Message: msg})
}

func NewHttpInternalServerError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, HttpError{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	})
}
