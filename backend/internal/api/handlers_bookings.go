package api

import (
	"net/http"
	"regexp"
	"time"

	"callcalendar/backend/internal/domain"
)

var emailRe = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

func (s *Server) createBooking(w http.ResponseWriter, r *http.Request) {
	var req BookingCreate
	if !decodeBody(w, r, &req) {
		return
	}
	if req.GuestName == "" || req.GuestEmail == "" {
		writeBadRequest(w, "guestName и guestEmail обязательны")
		return
	}
	if !emailRe.MatchString(req.GuestEmail) {
		writeBadRequest(w, "Некорректный email")
		return
	}
	startsAt, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		writeBadRequest(w, "startsAt должен быть в формате RFC3339")
		return
	}

	et, err := s.store.GetEventType(req.EventTypeID)
	if err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}

	now := s.now().UTC()
	if !domain.IsValidSlotStart(et, startsAt, now) {
		writeError(w, http.StatusUnprocessableEntity, "Слот недопустим")
		return
	}

	bk := domain.Booking{
		ID:          newID("bk"),
		EventTypeID: et.ID,
		StartsAt:    startsAt.UTC(),
		EndsAt:      startsAt.Add(time.Duration(et.DurationMinutes) * time.Minute).UTC(),
		GuestName:   req.GuestName,
		GuestEmail:  req.GuestEmail,
	}

	ok, err := s.store.CreateBooking(bk)
	if err != nil {
		writeInternal(w)
		return
	}
	if !ok {
		writeError(w, http.StatusConflict, "Выбранное время уже занято")
		return
	}
	writeJSON(w, http.StatusCreated, toBookingDTO(bk))
}

func (s *Server) adminUpcomingBookings(w http.ResponseWriter, r *http.Request) {
	bookings, err := s.store.UpcomingBookings(s.now().UTC())
	if err != nil {
		writeInternal(w)
		return
	}
	dto := make([]BookingDTO, 0, len(bookings))
	for _, bk := range bookings {
		dto = append(dto, toBookingDTO(bk))
	}
	writeJSON(w, http.StatusOK, dto)
}

func (s *Server) adminCancelBooking(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("bookingId")
	if err := s.store.DeleteBooking(id); err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toBookingDTO(bk domain.Booking) BookingDTO {
	return BookingDTO{
		ID:          bk.ID,
		EventTypeID: bk.EventTypeID,
		StartsAt:    bk.StartsAt.Format(time.RFC3339),
		EndsAt:      bk.EndsAt.Format(time.RFC3339),
		GuestName:   bk.GuestName,
		GuestEmail:  bk.GuestEmail,
	}
}
