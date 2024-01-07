package user

import (
	"database/sql"
	"errors"
	"net/http"
	"snippetbox/foundation/logger"
	"snippetbox/httperror"

	"github.com/alexedwards/scs/v2"
)

type Handler struct {
	errorHandler   *httperror.Handler
	sessionManager *scs.SessionManager
	UserService    *UserService
}

func NewHandler(logger *logger.Logger, sessionManager *scs.SessionManager, db *sql.DB) *Handler {
	return &Handler{
		errorHandler:   httperror.NewHandler(logger),
		UserService:    NewUserService(db),
		sessionManager: sessionManager,
	}
}

func (h *Handler) GetUserSignUpForm(w http.ResponseWriter, r *http.Request) {
	err := SignUpForm(r).Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering form")
		return
	}
}

func (h *Handler) GetUserLoginForm(w http.ResponseWriter, r *http.Request) {
	err := LoginForm(r).Render(r.Context(), w)

	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Error rendering form")
		return
	}
}

func (h *Handler) HandleSignUpUser(w http.ResponseWriter, r *http.Request) {

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	signupData := NewSignUpData(username, email, password)

	validationErrors := signupData.Validate()

	_ = validationErrors

	_, err := h.UserService.CreateUser(signupData)
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "internal server error")
	}

	// h.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *Handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	// Get the email and password from the form
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := h.UserService.Login(email, password)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			// Invalid credentials (e.g., wrong email or password)
			h.errorHandler.Unauthorized(w, r, "Invalid email or password. Please check your credentials and try again.")
			return
		case errors.Is(err, ErrUserNotFound):
			// User not found in the database
			h.errorHandler.NotFound(w, r, "User not found. Please sign up to create an account.")
			return
		default:
			// Any other unexpected error
			h.errorHandler.InternalServerError(w, r, err, "Oops! Something went wrong. Please try again later.")
			return
		}
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Oops! Something went wrong. Please try again later.")
		return
	}

	h.sessionManager.Put(r.Context(), "authenticatedUserID", user.ID)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *Handler) HandleLogoutUser(w http.ResponseWriter, r *http.Request) {
	err := h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.errorHandler.InternalServerError(w, r, err, "Oops! Something went wrong. Please try again later.")
		return
	}

	h.sessionManager.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}
