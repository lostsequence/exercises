package inmemory

import "movies-auth/users/internal/domain"

type UsersStorage struct {
	users []domain.User
}

func NewUsersStorage() *UsersStorage {
	return &UsersStorage{
		users: make([]domain.User, 0),
	}
}

func (s *UsersStorage) Insert(user domain.User) (domain.User, error) {
	var lastID int
	if len(s.users) > 0 {
		lastID = s.users[:len(s.users)-1][0].ID
	}

	user.ID = lastID + 1

	s.users = append(s.users, user)

	return user, nil
}

func (s *UsersStorage) IsUserExist(login string) (bool, error) {
	for i := range s.users {
		if s.users[i].Login == login {
			return true, domain.ErrNotFound
		}
	}

	return false, nil
}
