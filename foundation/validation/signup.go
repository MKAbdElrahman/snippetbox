package validation

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Validator interface defines the validation method
type Validator interface {
	Validate() ([]ValidationError, error)
}
