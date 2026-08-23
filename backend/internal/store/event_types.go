package store

import (
	"database/sql"
	"errors"

	"callcalendar/backend/internal/domain"
)

const eventTypeColumns = "id, name, description, duration_minutes, available_from, available_to"

func (s *Store) CreateEventType(et domain.EventType) error {
	_, err := s.db.Exec(
		`INSERT INTO event_types (id, name, description, duration_minutes, available_from, available_to)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		et.ID, et.Name, et.Description, et.DurationMinutes, et.AvailableFrom, et.AvailableTo,
	)
	return err
}

func (s *Store) ListEventTypes() ([]domain.EventType, error) {
	rows, err := s.db.Query(`SELECT ` + eventTypeColumns + ` FROM event_types ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []domain.EventType
	for rows.Next() {
		var et domain.EventType
		if err := rows.Scan(&et.ID, &et.Name, &et.Description, &et.DurationMinutes, &et.AvailableFrom, &et.AvailableTo); err != nil {
			return nil, err
		}
		types = append(types, et)
	}
	return types, rows.Err()
}

func (s *Store) GetEventType(id string) (domain.EventType, error) {
	row := s.db.QueryRow(`SELECT `+eventTypeColumns+` FROM event_types WHERE id = ?`, id)
	var et domain.EventType
	err := row.Scan(&et.ID, &et.Name, &et.Description, &et.DurationMinutes, &et.AvailableFrom, &et.AvailableTo)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.EventType{}, ErrNotFound
	}
	if err != nil {
		return domain.EventType{}, err
	}
	return et, nil
}

func (s *Store) UpdateEventType(et domain.EventType) error {
	res, err := s.db.Exec(
		`UPDATE event_types
		 SET name = ?, description = ?, duration_minutes = ?, available_from = ?, available_to = ?
		 WHERE id = ?`,
		et.Name, et.Description, et.DurationMinutes, et.AvailableFrom, et.AvailableTo, et.ID,
	)
	if err != nil {
		return err
	}
	return affectedAsErr(res)
}

func (s *Store) DeleteEventType(id string) error {
	res, err := s.db.Exec(`DELETE FROM event_types WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return affectedAsErr(res)
}

func affectedAsErr(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}
