package user

import (
	"regexp"
	"snippetbox/foundation/validation"
	"strings"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword string
	Created        time.Time
}

// SignUpData is a structure representing user sign-up data
type SignUpData struct {
	Username string
	Email    string
	Password string
}

// NewSignUpData creates a new SignUpData instance
func NewSignUpData(username, email, password string) *SignUpData {
	return &SignUpData{
		Username: strings.TrimSpace(username),
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
	}
}

// Validate validates the user sign-up data
func (s *SignUpData) Validate() []validation.ValidationError {
	var validationErrors []validation.ValidationError

	// Validate username
	if s.Username == "" {
		validationErrors = append(validationErrors, validation.ValidationError{Field: "username", Message: "Username is required"})
	}

	// Validate email format using a simple regular expression
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(s.Email) {
		validationErrors = append(validationErrors, validation.ValidationError{Field: "email", Message: "Invalid email format"})
	}

	// Validate password strength (at least 8 characters)
	if len(s.Password) < 8 {
		validationErrors = append(validationErrors, validation.ValidationError{Field: "password", Message: "Password must be at least 8 characters"})
	}

	return validationErrors
}

type LoginData struct {
	Username string
	Email    string
	Password string
}

func NewLoginData(email, password string) *LoginData {
	return &LoginData{
		Email:    strings.TrimSpace(email),
		Password: strings.TrimSpace(password),
	}
}

func (s *LoginData) Validate() []validation.ValidationError {
	var validationErrors []validation.ValidationError

	// Validate email format using a simple regular expression
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(s.Email) {
		validationErrors = append(validationErrors, validation.ValidationError{Field: "email", Message: "Invalid email format"})
	}

	// Validate password strength (at least 8 characters)
	if len(s.Password) < 8 {
		validationErrors = append(validationErrors, validation.ValidationError{Field: "password", Message: "Password must be at least 8 characters"})
	}

	return validationErrors
}
