package services

import (
	"fmt"
	"movies-auth/users/internal/domain"
	"time"

	"github.com/google/uuid"
)

type SessionsStorage interface {
	DeleteSessionByKey(key uuid.UUID) error
	GetSessionByKey(key uuid.UUID) (domain.Session, error)
	InsertSession(session domain.Session) (domain.Session, error)
}

type SessionsService struct {
	Storage SessionsStorage
}

func NewSessionService(storage SessionsStorage) *SessionsService {
	return &SessionsService{
		Storage: storage,
	}
}

func (s *SessionsService) CreateSession(userId int) (domain.Session, error) {
	session := domain.Session{
		Key:       uuid.New(),
		UserID:    userId,
		StartedAt: time.Now().UTC(),
	}

	newSession, err := s.Storage.InsertSession(session)
	if err != nil {
		return domain.Session{}, fmt.Errorf("failed to create user session: %w", err)
	}

	return newSession, nil
}

func (s *SessionsService) ValidateSession(key uuid.UUID) (domain.Session, error) {
	existingSession, err := s.Storage.GetSessionByKey(key)
	if err != nil {
		return domain.Session{}, err
	}

	return existingSession, nil
}

func (s *SessionsService) DeleteSession(key uuid.UUID) error {
	err := s.Storage.DeleteSessionByKey(key)
	if err != nil {
		return err
	}

	return nil
}
