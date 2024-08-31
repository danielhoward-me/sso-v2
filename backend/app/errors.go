package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Status  int    `json:"status" binding:"required"`
	Error   string `json:"error" binding:"required"`
	Message string `json:"message"`
}

var errorMessageHandlers = map[int](func(*gin.Context) string){
	http.StatusNotFound: func(c *gin.Context) string {
		return fmt.Sprintf("%s %s is not a valid API endpoint", c.Request.Method, c.Request.URL.Path)
	},
	http.StatusInternalServerError: func(c *gin.Context) string {
		return "There was an unkown internal server error when processing your request"
	},
}

func handleError(code int, c *gin.Context) {
	errorBody := errorResponse{
		Status: code,
		Error:  http.StatusText(code),
	}

	handler, exists := errorMessageHandlers[code]
	if exists {
		errorBody.Message = handler(c)
	} else {
		fmt.Printf("No handler could be found for http code %d, returning without a message\n", code)
	}

	c.JSON(code, errorBody)
}
