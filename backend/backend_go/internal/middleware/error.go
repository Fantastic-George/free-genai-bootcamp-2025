package middleware

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// ErrorHandler middleware handles errors in a standardized way
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Only handle errors if there are any
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			var response ErrorResponse

			switch {
			case errors.Is(err, sql.ErrNoRows):
				response = ErrorResponse{
					Status:  http.StatusNotFound,
					Message: "Resource not found",
					Details: err.Error(),
				}

			case errors.As(err, &validator.ValidationErrors{}):
				response = ErrorResponse{
					Status:  http.StatusBadRequest,
					Message: "Validation error",
					Details: err.Error(),
				}

			default:
				// Log unexpected errors
				log.Printf("Unexpected error: %v", err)
				response = ErrorResponse{
					Status:  http.StatusInternalServerError,
					Message: "Internal server error",
					Details: "An unexpected error occurred",
				}
			}

			c.JSON(response.Status, response)
			c.Abort()
		}
	}
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateRequest validates a request struct and returns formatted errors
func ValidateRequest(obj interface{}) []ValidationError {
	validate := validator.New()
	err := validate.Struct(obj)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Message: getValidationErrorMsg(err),
		})
	}

	return validationErrors
}

// getValidationErrorMsg returns a human-readable validation error message
func getValidationErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too small"
	case "max":
		return "Value is too large"
	case "email":
		return "Invalid email format"
	default:
		return "Invalid value"
	}
}
