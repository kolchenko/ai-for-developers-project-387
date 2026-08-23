package api

type EventTypeDTO struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
	AvailableFrom   string `json:"availableFrom"`
	AvailableTo     string `json:"availableTo"`
}

type EventTypeCreate struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
	AvailableFrom   string `json:"availableFrom"`
	AvailableTo     string `json:"availableTo"`
}

type EventTypeUpdate struct {
	Name            *string `json:"name"`
	Description     *string `json:"description"`
	DurationMinutes *int    `json:"durationMinutes"`
	AvailableFrom   *string `json:"availableFrom"`
	AvailableTo     *string `json:"availableTo"`
}

type BookingDTO struct {
	ID          string `json:"id"`
	EventTypeID string `json:"eventTypeId"`
	StartsAt    string `json:"startsAt"`
	EndsAt      string `json:"endsAt"`
	GuestName   string `json:"guestName"`
	GuestEmail  string `json:"guestEmail"`
}

type BookingCreate struct {
	EventTypeID string `json:"eventTypeId"`
	StartsAt    string `json:"startsAt"`
	GuestName   string `json:"guestName"`
	GuestEmail  string `json:"guestEmail"`
}

type SlotDTO struct {
	StartsAt string `json:"startsAt"`
	EndsAt   string `json:"endsAt"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Username string `json:"username"`
}
