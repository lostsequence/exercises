package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"movies-auth/users/internal/api/middlewares"
	"movies-auth/users/internal/domain"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UsersService interface {
	Login(ctx context.Context, user domain.User) (domain.User, error)
	Create(ctx context.Context, user domain.User) (domain.User, error)
}

type SessionService interface {
	DeleteSession(key uuid.UUID) error
	CreateSession(userId int) (domain.Session, error)
	ValidateSession(key uuid.UUID) (domain.Session, error)
}

type UsersHandler struct {
	UsersService    UsersService
	SessionsService SessionService
}

func NewUsersHandler(usersService UsersService, sessionsService SessionService) UsersHandler {
	return UsersHandler{
		UsersService:    usersService,
		SessionsService: sessionsService,
	}
}

//go:generate mockgen -source users.go -destination ../../tests/api_mocks/users.go package apimocks

func (h UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "not supported content type", http.StatusUnsupportedMediaType)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	createdUser, err := h.UsersService.Create(r.Context(), user)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrConflict):
			http.Error(w, "user already exist", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}

		log.Println(err)
		return
	}

	userSession, err := h.SessionsService.CreateSession(createdUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "session",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 5),
		Value:    userSession.Key.String(),
	}

	userBytes, err := json.Marshal(createdUser)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, sessionCookie)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(userBytes)
}

func (h UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "not supported content type", http.StatusUnsupportedMediaType)
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	loggedUser, err := h.UsersService.Login(r.Context(), user)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusBadRequest)
		return
	}

	userSession, err := h.SessionsService.CreateSession(loggedUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	sessionCookie := &http.Cookie{
		Name:     "session",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Minute * 5),
		Value:    userSession.Key.String(),
	}

	http.SetCookie(w, sessionCookie)
	w.WriteHeader(http.StatusOK)
}

func (h UsersHandler) Logout(w http.ResponseWriter, r *http.Request) {
	sessKey := r.Context().Value(middlewares.SessionKey).(uuid.UUID)
	w.Write([]byte(sessKey.String()))
	w.WriteHeader(http.StatusOK)

	log.Println(sessKey.String())
	err := h.SessionsService.DeleteSession(sessKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "unexpected error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h UsersHandler) Session(w http.ResponseWriter, r *http.Request) {
	sessionKeyParam := chi.URLParam(r, "key")
	if sessionKeyParam == "" {
		log.Println("session key required")
		http.Error(w, "session key required", http.StatusBadRequest)
		return
	}

	sessionKey, err := uuid.Parse(sessionKeyParam)
	if err != nil {
		log.Println(err)
		http.Error(w, "invalid session key", http.StatusInternalServerError)
		return
	}

	session, err := h.SessionsService.ValidateSession(sessionKey)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionBytes, err := json.Marshal(session)
	if err != nil {
		log.Println(err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(sessionBytes)
}

func (h UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("param1")
	log.Println(param)

	w.Write([]byte("list of users"))
	w.WriteHeader(http.StatusOK)
}
