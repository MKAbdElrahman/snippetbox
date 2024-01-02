package user

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		DB: db,
	}
}

func (s *UserService) CreateUser(signupData *SignUpData) (User, error) {

	hashed_password, err := bcrypt.GenerateFromPassword([]byte(signupData.Password), 12)
	if err != nil {
		return User{}, fmt.Errorf("failed to generate password hash: %w", err)
	}

	result, err := s.DB.Exec("INSERT INTO users (name, email, hashed_password, created) VALUES (?, ?, ?, ?)",
		signupData.Username, signupData.Email, hashed_password, time.Now())
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return User{
		ID:             int(userID),
		Name:           signupData.Username,
		Email:          signupData.Email,
		HashedPassword: string(hashed_password),
		Created:        time.Now(),
	}, nil
}

func (s *UserService) Login(email, password string) (User, error) {
	var user User

	// Fetch the user by email
	err := s.DB.QueryRow("SELECT id, name, email, hashed_password, created FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Created)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Compare the provided password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return User{}, ErrInvalidCredentials
	}

	return user, nil
}
