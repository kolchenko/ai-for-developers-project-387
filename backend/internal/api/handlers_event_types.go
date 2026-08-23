package api

import (
	"crypto/rand"
	"net/http"

	"callcalendar/backend/internal/domain"
	"callcalendar/backend/internal/store"
)

func (s *Server) listEventTypes(w http.ResponseWriter, r *http.Request) {
	types, err := s.store.ListEventTypes()
	if err != nil {
		writeInternal(w)
		return
	}
	dto := make([]EventTypeDTO, 0, len(types))
	for _, et := range types {
		dto = append(dto, toEventTypeDTO(et))
	}
	writeJSON(w, http.StatusOK, dto)
}

func (s *Server) adminCreateEventType(w http.ResponseWriter, r *http.Request) {
	var req EventTypeCreate
	if !decodeBody(w, r, &req) {
		return
	}
	et, msg := validateEventType(req.Name, req.Description, req.DurationMinutes, req.AvailableFrom, req.AvailableTo)
	if msg != "" {
		writeBadRequest(w, msg)
		return
	}
	et.ID = newID("et")
	if err := s.store.CreateEventType(et); err != nil {
		writeInternal(w)
		return
	}
	writeJSON(w, http.StatusCreated, toEventTypeDTO(et))
}

func (s *Server) adminUpdateEventType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("eventTypeId")
	current, err := s.store.GetEventType(id)
	if err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}

	var patch EventTypeUpdate
	if !decodeBody(w, r, &patch) {
		return
	}

	next := current
	if patch.Name != nil {
		next.Name = *patch.Name
	}
	if patch.Description != nil {
		next.Description = *patch.Description
	}
	if patch.DurationMinutes != nil {
		next.DurationMinutes = *patch.DurationMinutes
	}
	if patch.AvailableFrom != nil {
		next.AvailableFrom = *patch.AvailableFrom
	}
	if patch.AvailableTo != nil {
		next.AvailableTo = *patch.AvailableTo
	}

	if _, msg := validateEventType(next.Name, next.Description, next.DurationMinutes, next.AvailableFrom, next.AvailableTo); msg != "" {
		writeBadRequest(w, msg)
		return
	}
	if err := s.store.UpdateEventType(next); err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toEventTypeDTO(next))
}

func (s *Server) adminDeleteEventType(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("eventTypeId")
	if err := s.store.DeleteEventType(id); err != nil {
		writeNotFoundOrInternal(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateEventType(
	name, description string,
	durationMinutes int,
	availableFrom, availableTo string,
) (domain.EventType, string) {
	if name == "" {
		return domain.EventType{}, "name обязателен"
	}
	if description == "" {
		return domain.EventType{}, "description обязательна"
	}
	if !domain.IsValidDuration(durationMinutes) {
		return domain.EventType{}, "durationMinutes должен быть одним из: 15, 30, 45, 60"
	}
	from, errFrom := domain.ParseTimeOfDay(availableFrom)
	to, errTo := domain.ParseTimeOfDay(availableTo)
	if errFrom != nil || errTo != nil {
		return domain.EventType{}, "availableFrom и availableTo должны быть в формате HH:mm:ss"
	}
	if !(from < to) {
		return domain.EventType{}, "availableFrom должен быть раньше availableTo"
	}
	return domain.EventType{
		Name:            name,
		Description:     description,
		DurationMinutes: durationMinutes,
		AvailableFrom:   availableFrom,
		AvailableTo:     availableTo,
	}, ""
}

func toEventTypeDTO(et domain.EventType) EventTypeDTO {
	return EventTypeDTO{
		ID:              et.ID,
		Name:            et.Name,
		Description:     et.Description,
		DurationMinutes: et.DurationMinutes,
		AvailableFrom:   et.AvailableFrom,
		AvailableTo:     et.AvailableTo,
	}
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) bool {
	if err := jsonDecode(r, dst); err != nil {
		writeBadRequest(w, "Некорректное тело запроса")
		return false
	}
	return true
}

func writeNotFoundOrInternal(w http.ResponseWriter, err error) {
	if err == store.ErrNotFound {
		writeNotFound(w)
		return
	}
	writeInternal(w)
}

func newID(prefix string) string {
	const alphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return prefix + "-" + string(b)
}
