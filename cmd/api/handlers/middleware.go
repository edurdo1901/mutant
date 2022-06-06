package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"prueba.com/internal/mutant"
)

// errorHandler process errors generated in the handler.
func errorHandler(c *gin.Context) {
	c.Next()
	if len(c.Errors) > 0 {
		for _, err := range c.Errors {
			code := getStatusError(err)
			c.JSON(code, Error{
				Code:    http.StatusText(code),
				Message: err.Error(),
			})
		}
	}
}

// getStatusError get the status code according to error.
func getStatusError(err error) int {
	var statusCode int
	var validationErr validator.ValidationErrors
	switch {
	case errors.Is(err, mutant.ErrInvalidDna),
		errors.Is(err, mutant.ErrInvalidLength),
		errors.As(err, &validationErr):
		statusCode = http.StatusUnprocessableEntity
	default:
		statusCode = http.StatusInternalServerError
	}

	return statusCode
}
