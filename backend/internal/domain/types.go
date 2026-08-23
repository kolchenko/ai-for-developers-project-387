package domain

import "time"

type EventType struct {
	ID              string
	Name            string
	Description     string
	DurationMinutes int
	AvailableFrom   string
	AvailableTo     string
}

type Booking struct {
	ID          string
	EventTypeID string
	StartsAt    time.Time
	EndsAt      time.Time
	GuestName   string
	GuestEmail  string
}

type Slot struct {
	StartsAt time.Time
	EndsAt   time.Time
}

type Admin struct {
	Username     string
	PasswordHash string
}
