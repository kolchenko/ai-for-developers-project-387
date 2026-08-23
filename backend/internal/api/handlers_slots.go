package api

import (
	"net/http"
	"time"

	"callcalendar/backend/internal/domain"
)

func (s *Server) getSlots(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("eventTypeId")
	et, err := s.store.GetEventType(id)
	if err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}

	now := s.now().UTC()
	y, m, d := now.Date()
	start := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	end := start.Add(domain.BookingWindowDays * 24 * time.Hour)

	bookings, err := s.store.BookingsOverlapping(start, end)
	if err != nil {
		writeInternal(w)
		return
	}

	slots := domain.GenerateGridSlots(et, now)
	free := domain.FilterFree(slots, bookings)

	dto := make([]SlotDTO, 0, len(free))
	for _, slot := range free {
		dto = append(dto, SlotDTO{
			StartsAt: slot.StartsAt.Format(time.RFC3339),
			EndsAt:   slot.EndsAt.Format(time.RFC3339),
		})
	}
	writeJSON(w, http.StatusOK, dto)
}
