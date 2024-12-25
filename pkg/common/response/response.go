package response

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   any         `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, message string, data interface{}) {
	c.JSON(status, Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, status int, message string, err any) {
	c.AbortWithStatusJSON(status, Response{
		Status:  status,
		Message: message,
		Error:   err,
	})
}

func RenderBindingErrors(ctx *gin.Context, validationError validator.ValidationErrors) {
	var responseErrs []any
	for _, fieldError := range validationError {
		field := fieldError.Field()
		f := strings.ToLower(field[:1]) + field[1:]
		responseErrs = append(responseErrs, gin.H{
			f: fieldError.ActualTag(),
		})
	}
	Error(ctx, http.StatusBadRequest, "Invalid input", responseErrs)
}
