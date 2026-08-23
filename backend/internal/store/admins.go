package store

import (
	"database/sql"
	"errors"

	"callcalendar/backend/internal/domain"
)

func (s *Store) GetAdmin(username string) (domain.Admin, error) {
	row := s.db.QueryRow(`SELECT username, password_hash FROM admins WHERE username = ?`, username)
	var admin domain.Admin
	err := row.Scan(&admin.Username, &admin.PasswordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Admin{}, ErrNotFound
	}
	if err != nil {
		return domain.Admin{}, err
	}
	return admin, nil
}
