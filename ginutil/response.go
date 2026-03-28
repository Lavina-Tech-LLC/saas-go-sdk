package ginutil

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard JSON envelope for Gin handlers.
type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	IsOk    bool        `json:"isOk"`
}

// OK sends a 200 JSON response with isOk: true.
func OK(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusOK, Response{Data: data, Message: msg, IsOk: true})
}

// Created sends a 201 JSON response with isOk: true.
func Created(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	c.JSON(http.StatusCreated, Response{Data: data, Message: msg, IsOk: true})
}

// Error sends an error JSON response with isOk: false and the given HTTP status code.
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{Data: nil, Message: message, IsOk: false})
}

// Fail sends a 400 Bad Request JSON response with isOk: false.
func Fail(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

// Unauthorized sends a 401 JSON response with isOk: false.
func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message)
}

// NotFound sends a 404 JSON response with isOk: false.
func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}

// InternalError sends a 500 JSON response with isOk: false.
func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, message)
}
