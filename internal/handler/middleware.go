package handler

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// APIKeyAuth blocks requests whose X-API-Key header does not match
// env variable BUDGET_API_KEY. If the env var is unset, startup aborts.
func APIKeyAuth() gin.HandlerFunc {
	expected := os.Getenv("BUDGET_API_KEY")
	if expected == "" {
		panic("BUDGET_API_KEY not set")
	}
	return func(c *gin.Context) {
		if c.GetHeader("X-API-Key") != expected {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "invalid API key"})
			return
		}
		c.Next()
	}
}

// ValidateRequest is a middleware that validates request body against a struct
// using validator v10. It expects the struct to be passed as a type parameter.
func ValidateRequest[T any]() gin.HandlerFunc {
	validate := validator.New()
	
	// Register custom validators if needed
	registerCustomValidators(validate)
	
	return func(c *gin.Context) {
		var request T
		
		// Bind JSON to struct
		if err := c.ShouldBindJSON(&request); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid request format",
				"data":  nil,
			})
			return
		}
		
		// Validate struct
		if err := validate.Struct(request); err != nil {
			validationErrors := make(map[string]string)
			
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, fieldError := range ve {
					field := fieldError.Field()
					tag := fieldError.Tag()
					param := fieldError.Param()
					
					// Convert field name to snake_case for API consistency
					fieldName := toSnakeCase(field)
					
					// Create user-friendly error messages
					message := getValidationMessage(tag, param)
					validationErrors[fieldName] = message
				}
			}
			
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "validation failed",
				"data":  validationErrors,
			})
			return
		}
		
		// Store validated request in context for handlers to use
		c.Set("validated_request", request)
		c.Next()
	}
}

// registerCustomValidators registers any custom validation functions
func registerCustomValidators(v *validator.Validate) {
	// Register currency validator for amount fields
	v.RegisterValidation("currency", validateCurrency)
	// Register date validator for date fields
	v.RegisterValidation("date", validateDate)
}

// validateCurrency validates currency format (e.g., "-12.34", "123.45")
func validateCurrency(fl validator.FieldLevel) bool {
	amount := fl.Field().String()
	
	// Check if empty (handled by required validator)
	if amount == "" {
		return true
	}
	
	// Remove leading minus sign if present
	cleanAmount := amount
	if strings.HasPrefix(amount, "-") {
		cleanAmount = amount[1:]
	}
	
	// Check if it's a valid decimal number
	parts := strings.Split(cleanAmount, ".")
	if len(parts) != 2 {
		return false
	}
	
	// Validate integer part
	if parts[0] == "" {
		return false
	}
	
	// Validate decimal part (must be exactly 2 digits and numeric)
	if len(parts[1]) != 2 {
		return false
	}
	
	for _, r := range parts[1] {
		if r < '0' || r > '9' {
			return false
		}
	}
	
	// Try to parse as float to ensure it's a valid number
	_, err := strconv.ParseFloat(amount, 64)
	return err == nil
}

// validateDate validates date format (YYYY-MM-DD)
func validateDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	
	// Check if empty (handled by required validator)
	if dateStr == "" {
		return true
	}
	
	// Try to parse the date in YYYY-MM-DD format
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// toSnakeCase converts camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// getValidationMessage returns user-friendly validation error messages
func getValidationMessage(tag, param string) string {
	switch tag {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return "must be at least " + param + " characters"
	case "max":
		return "must be at most " + param + " characters"
	case "oneof":
		return "must be one of: " + param
	case "gte":
		return "must be greater than or equal to " + param
	case "lte":
		return "must be less than or equal to " + param
	case "gt":
		return "must be greater than " + param
	case "lt":
		return "must be less than " + param
	case "date":
		return "must be a valid date in YYYY-MM-DD format"
	case "datetime":
		return "must be a valid datetime"
	case "url":
		return "must be a valid URL"
	case "numeric":
		return "must be a valid number"
	case "alpha":
		return "must contain only letters"
	case "alphanum":
		return "must contain only letters and numbers"
	case "currency":
		return "must be a valid currency amount (e.g., '12.34' or '-12.34')"
	default:
		return "validation failed for " + tag
	}
}

// GetValidatedRequest extracts the validated request from the Gin context
func GetValidatedRequest[T any](c *gin.Context) (T, bool) {
	value, exists := c.Get("validated_request")
	if !exists {
		var zero T
		return zero, false
	}
	
	request, ok := value.(T)
	return request, ok
} 