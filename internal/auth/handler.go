package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service  Service
	validate *validator.Validate
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service:  service,
		validate: validator.New(),
	}
}

type errorResponse struct {
	Error string `json:"error"`
}

type registerResponse struct {
	UserId UserId `json:"user_id"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}

func RegisterHandler(h *Handler) {
	http.HandleFunc("/auth/register", h.RegisterHandler)
	http.HandleFunc("/auth/login", h.LoginHandler)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	userId, err := h.service.RegisterUser(r.Context(), input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrUserWithEmailAlreadyExists) || errors.Is(err, ErrUserWithUsernameAlreadyExists) {
			statusCode = http.StatusConflict
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponse{UserId: userId})
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var input LoginUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: "invalid request body"})
		return
	}

	if err := h.validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	response, err := h.service.LoginUser(r.Context(), input)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, ErrInvalidCredentials) {
			statusCode = http.StatusUnauthorized
		}
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errorResponse{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
