package store

import (
	"database/sql"
	"errors"
	"time"

	"callcalendar/backend/internal/domain"
)

const bookingColumns = "id, event_type_id, starts_at, ends_at, guest_name, guest_email"

// CreateBooking атомарно проверяет пересечение с существующими бронями и вставляет
// новую. Возвращает ok=false, если интервал занят (409).
func (s *Store) CreateBooking(bk domain.Booking) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	overlap, err := overlaps(tx, bk.StartsAt, bk.EndsAt)
	if err != nil {
		return false, err
	}
	if overlap {
		return false, nil
	}

	_, err = tx.Exec(
		`INSERT INTO bookings (id, event_type_id, starts_at, ends_at, guest_name, guest_email)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		bk.ID, bk.EventTypeID, formatTime(bk.StartsAt), formatTime(bk.EndsAt),
		bk.GuestName, bk.GuestEmail,
	)
	if err != nil {
		return false, err
	}
	return true, tx.Commit()
}

// BookingsOverlapping возвращает все брони, пересекающие интервал [start, end).
func (s *Store) BookingsOverlapping(start, end time.Time) ([]domain.Booking, error) {
	rows, err := s.db.Query(
		`SELECT `+bookingColumns+` FROM bookings WHERE starts_at < ? AND ends_at > ?`,
		formatTime(end), formatTime(start),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		bk, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, bk)
	}
	return bookings, rows.Err()
}

func (s *Store) UpcomingBookings(now time.Time) ([]domain.Booking, error) {
	rows, err := s.db.Query(
		`SELECT `+bookingColumns+` FROM bookings WHERE starts_at >= ? ORDER BY starts_at`,
		formatTime(now),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []domain.Booking
	for rows.Next() {
		bk, err := scanBooking(rows)
		if err != nil {
			return nil, err
		}
		bookings = append(bookings, bk)
	}
	return bookings, rows.Err()
}

func (s *Store) GetBooking(id string) (domain.Booking, error) {
	row := s.db.QueryRow(`SELECT `+bookingColumns+` FROM bookings WHERE id = ?`, id)
	bk, err := scanBooking(row)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Booking{}, ErrNotFound
	}
	if err != nil {
		return domain.Booking{}, err
	}
	return bk, nil
}

func (s *Store) DeleteBooking(id string) error {
	res, err := s.db.Exec(`DELETE FROM bookings WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return affectedAsErr(res)
}

type scanner interface {
	Scan(dest ...any) error
}

func scanBooking(row scanner) (domain.Booking, error) {
	var bk domain.Booking
	var startsAt, endsAt string
	if err := row.Scan(&bk.ID, &bk.EventTypeID, &startsAt, &endsAt, &bk.GuestName, &bk.GuestEmail); err != nil {
		return domain.Booking{}, err
	}
	st, err := time.Parse(timeFormat, startsAt)
	if err != nil {
		return domain.Booking{}, err
	}
	en, err := time.Parse(timeFormat, endsAt)
	if err != nil {
		return domain.Booking{}, err
	}
	bk.StartsAt, bk.EndsAt = st, en
	return bk, nil
}

func overlaps(q interface {
	QueryRow(query string, args ...any) *sql.Row
}, start, end time.Time) (bool, error) {
	var one int
	err := q.QueryRow(
		`SELECT 1 FROM bookings WHERE starts_at < ? AND ends_at > ? LIMIT 1`,
		formatTime(end), formatTime(start),
	).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
