package db

import (
	"context"
	"database/sql"
	"movies-auth/users/internal/domain"
)

type DbStorage struct {
	db *sql.DB
}

func NewDbStorage(dbCon *sql.DB) *DbStorage {
	return &DbStorage{
		db: dbCon,
	}
}

func (s *DbStorage) Insert(ctx context.Context, user domain.User) (domain.User, error) {
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id, login, password`

	var newUser domain.User
	err := s.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(&newUser.ID, &newUser.Login, &newUser.Password)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (s *DbStorage) GetUserByID(ctx context.Context, login string) (domain.User, error) {
	query := `SELECT id, login, password, notification_sent FROM users WHERE login = $1`

	var newUser domain.User
	err := s.db.QueryRowContext(ctx, query, login).Scan(&newUser.ID, &newUser.Login, &newUser.Password, &newUser.NotificationSent)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (s *DbStorage) IsUserExist(ctx context.Context, login string) (bool, error) {
	query := `SELECT id FROM users WHERE login = $1`
	var userID int
	err := s.db.QueryRowContext(ctx, query, login).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, domain.ErrNotFound
		}

		return false, err
	}

	return true, nil
}

func (s *DbStorage) UpdateNotificationSent(ctx context.Context, id int) error {
	query := `UPDATE users SET notification_sent = true WHERE id = $1`
	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *DbStorage) GetUsersWithExpiredPassword(ctx context.Context) ([]domain.User, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT * FROM users WHERE current_timestamp > password_expires AND notification_sent = false")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.Login, &user.Password, &user.PasswordExpires, &user.NotificationSent); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
