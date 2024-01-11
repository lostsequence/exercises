package db

import (
	"movies-auth/users/internal/domain"

	"github.com/google/uuid"
)

func (s *DbStorage) InsertSession(session domain.Session) (domain.Session, error) {
	query := `INSERT INTO sessions (key, userid, startedat) VALUES ($1, $2, $3) RETURNING id, key, userid, startedat`

	var newSession domain.Session
	err := s.db.
		QueryRow(query, session.Key, session.UserID, session.StartedAt).
		Scan(&newSession.ID, &newSession.Key, &newSession.UserID, &newSession.StartedAt)
	if err != nil {
		return domain.Session{}, err
	}

	return newSession, nil
}

func (s *DbStorage) GetSessionByKey(key uuid.UUID) (domain.Session, error) {
	query := `SELECT id, key, userid, startedat FROM sessions WHERE key = $1`

	var newSession domain.Session
	err := s.db.QueryRow(query, key).
		Scan(&newSession.ID, &newSession.Key, &newSession.UserID, &newSession.StartedAt)
	if err != nil {
		return domain.Session{}, err
	}

	return newSession, nil
}

func (s *DbStorage) DeleteSessionByKey(key uuid.UUID) error {
	query := `DELETE FROM sessions WHERE key = $1`

	_, err := s.db.Exec(query, key)
	if err != nil {
		return err
	}

	return nil
}
