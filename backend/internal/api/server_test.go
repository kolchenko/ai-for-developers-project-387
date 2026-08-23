package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"callcalendar/backend/internal/domain"
	"callcalendar/backend/internal/store"
)

const testNow = "2026-08-23T00:00:00Z"

var fixedNow = func() time.Time {
	t, _ := time.Parse(time.RFC3339, testNow)
	return t
}

type testEnv struct {
	ts    *httptest.Server
	store *store.Store
}

func TestPrefixedAPIRoutes(t *testing.T) {
	ts := newTestServer(t)

	status, body := doJSON(t, http.MethodGet, ts.URL+"/api/event-types", "")
	if status != http.StatusOK {
		t.Fatalf("prefixed list: status %d", status)
	}
	if got := decode[[]EventTypeDTO](t, body); len(got) != 0 {
		t.Fatalf("expected empty list, got %+v", got)
	}

	created := createEventType(t, ts, consultType)
	status, body = doJSON(t, http.MethodGet, ts.URL+"/api/event-types/"+created.ID+"/slots", "")
	if status != http.StatusOK {
		t.Fatalf("prefixed slots: status %d body %s", status, body)
	}
	if slots := decode[[]SlotDTO](t, body); len(slots) != 252 {
		t.Fatalf("expected 252 slots, got %d", len(slots))
	}
}

func TestSPAServing(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>spa</html>"), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	ts := httptest.NewServer(NewServerWithWeb(st, dir))
	t.Cleanup(ts.Close)

	res, err := http.Get(ts.URL + "/event-types/et-1")
	if err != nil {
		t.Fatalf("get spa route: %v", err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK || string(body) != "<html>spa</html>" {
		t.Fatalf("spa fallback: status %d body %s", res.StatusCode, body)
	}
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { st.Close() })

	ts := httptest.NewServer(NewTestServer(st, fixedNow))
	t.Cleanup(ts.Close)
	return &testEnv{ts: ts, store: st}
}

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	return newTestEnv(t).ts
}

func doJSON(t *testing.T, method, url, body string) (int, []byte) {
	t.Helper()
	var req *http.Request
	var err error
	if body == "" {
		req, err = http.NewRequest(method, url, nil)
	} else {
		req, err = http.NewRequest(method, url, bytes.NewBufferString(body))
	}
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("%s %s: %v", method, url, err)
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(res.Body); err != nil {
		t.Fatalf("read body: %v", err)
	}
	return res.StatusCode, buf.Bytes()
}

func decode[T any](t *testing.T, body []byte) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatalf("decode body %q: %v", body, err)
	}
	return v
}

func createEventType(t *testing.T, ts *httptest.Server, payload string) EventTypeDTO {
	t.Helper()
	status, body := doJSON(t, http.MethodPost, ts.URL+"/admin/event-types", payload)
	if status != http.StatusCreated {
		t.Fatalf("create event type: status %d body %s", status, body)
	}
	return decode[EventTypeDTO](t, body)
}

const consultType = `{
	"name": "Консультация",
	"description": "Консультация по вашему проекту",
	"durationMinutes": 30,
	"availableFrom": "09:00:00",
	"availableTo": "18:00:00"
}`

const reviewType = `{
	"name": "Разбор кода",
	"description": "Код-ревью",
	"durationMinutes": 45,
	"availableFrom": "10:00:00",
	"availableTo": "17:00:00"
}`

func TestEventTypesCRUD(t *testing.T) {
	ts := newTestServer(t)

	status, body := doJSON(t, http.MethodGet, ts.URL+"/event-types", "")
	if status != http.StatusOK {
		t.Fatalf("list: status %d", status)
	}
	if got := decode[[]EventTypeDTO](t, body); len(got) != 0 {
		t.Fatalf("expected empty list, got %+v", got)
	}

	created := createEventType(t, ts, consultType)
	if created.ID == "" || created.Name != "Консультация" || created.DurationMinutes != 30 {
		t.Fatalf("unexpected created: %+v", created)
	}

	status, body = doJSON(t, http.MethodGet, ts.URL+"/event-types", "")
	if got := decode[[]EventTypeDTO](t, body); len(got) != 1 || got[0].ID != created.ID {
		t.Fatalf("unexpected list after create: %+v", got)
	}

	// PATCH: частичное обновление
	patch := `{"name":"Консультация 30 мин","availableTo":"19:00:00"}`
	status, body = doJSON(t, http.MethodPatch, ts.URL+"/admin/event-types/"+created.ID, patch)
	if status != http.StatusOK {
		t.Fatalf("patch: status %d body %s", status, body)
	}
	updated := decode[EventTypeDTO](t, body)
	if updated.Name != "Консультация 30 мин" || updated.AvailableTo != "19:00:00" || updated.DurationMinutes != 30 {
		t.Fatalf("unexpected updated: %+v", updated)
	}

	// PATCH несуществующего
	status, _ = doJSON(t, http.MethodPatch, ts.URL+"/admin/event-types/et-nope", patch)
	if status != http.StatusNotFound {
		t.Fatalf("patch missing: status %d, want 404", status)
	}

	// DELETE
	status, _ = doJSON(t, http.MethodDelete, ts.URL+"/admin/event-types/"+created.ID, "")
	if status != http.StatusNoContent {
		t.Fatalf("delete: status %d, want 204", status)
	}
	status, body = doJSON(t, http.MethodGet, ts.URL+"/event-types", "")
	if got := decode[[]EventTypeDTO](t, body); status != http.StatusOK || len(got) != 0 {
		t.Fatalf("list after delete: status %d len %d body %s", status, len(got), body)
	}
	status, _ = doJSON(t, http.MethodDelete, ts.URL+"/admin/event-types/"+created.ID, "")
	if status != http.StatusNotFound {
		t.Fatalf("delete missing: status %d, want 404", status)
	}
}

func TestCreateEventTypeValidation(t *testing.T) {
	ts := newTestServer(t)

	cases := []struct {
		name    string
		payload string
	}{
		{"empty name", `{"name":"","description":"d","durationMinutes":30,"availableFrom":"09:00:00","availableTo":"18:00:00"}`},
		{"bad duration", `{"name":"n","description":"d","durationMinutes":20,"availableFrom":"09:00:00","availableTo":"18:00:00"}`},
		{"from after to", `{"name":"n","description":"d","durationMinutes":30,"availableFrom":"18:00:00","availableTo":"09:00:00"}`},
		{"bad time format", `{"name":"n","description":"d","durationMinutes":30,"availableFrom":"9:00","availableTo":"18:00:00"}`},
		{"malformed json", `{`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, _ := doJSON(t, http.MethodPost, ts.URL+"/admin/event-types", tc.payload)
			if status != http.StatusBadRequest {
				t.Fatalf("status %d, want 400", status)
			}
		})
	}
}

func TestSlots(t *testing.T) {
	ts := newTestServer(t)
	created := createEventType(t, ts, consultType)

	// несуществующий тип
	status, _ := doJSON(t, http.MethodGet, ts.URL+"/event-types/et-nope/slots", "")
	if status != http.StatusNotFound {
		t.Fatalf("slots missing type: status %d, want 404", status)
	}

	status, body := doJSON(t, http.MethodGet, ts.URL+"/event-types/"+created.ID+"/slots", "")
	if status != http.StatusOK {
		t.Fatalf("slots: status %d body %s", status, body)
	}
	slots := decode[[]SlotDTO](t, body)
	if len(slots) != 14*18 {
		t.Fatalf("expected 252 free slots, got %d", len(slots))
	}
	if slots[0].StartsAt != "2026-08-23T09:00:00Z" {
		t.Fatalf("unexpected first slot: %s", slots[0].StartsAt)
	}

	// бронь на 09:00 — слот должен исчезнуть из списка
	status, body = doJSON(t, http.MethodPost, ts.URL+"/bookings",
		`{"eventTypeId":"`+created.ID+`","startsAt":"2026-08-24T09:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`)
	if status != http.StatusCreated {
		t.Fatalf("book: status %d body %s", status, body)
	}

	status, body = doJSON(t, http.MethodGet, ts.URL+"/event-types/"+created.ID+"/slots", "")
	slots = decode[[]SlotDTO](t, body)
	if len(slots) != 251 {
		t.Fatalf("expected 251 free slots after booking, got %d", len(slots))
	}
	for _, s := range slots {
		if s.StartsAt == "2026-08-24T09:00:00Z" {
			t.Fatal("booked slot still returned")
		}
	}
}

func TestCreateBooking(t *testing.T) {
	ts := newTestServer(t)
	consult := createEventType(t, ts, consultType)

	// успешное бронирование
	payload := `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-24T09:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`
	status, body := doJSON(t, http.MethodPost, ts.URL+"/bookings", payload)
	if status != http.StatusCreated {
		t.Fatalf("book: status %d body %s", status, body)
	}
	bk := decode[BookingDTO](t, body)
	if bk.ID == "" || bk.EventTypeID != consult.ID || bk.StartsAt != "2026-08-24T09:00:00Z" ||
		bk.EndsAt != "2026-08-24T09:30:00Z" || bk.GuestName != "Иван" {
		t.Fatalf("unexpected booking: %+v", bk)
	}

	cases := []struct {
		name    string
		payload string
		want    int
	}{
		{"not found type", `{"eventTypeId":"et-nope","startsAt":"2026-08-24T09:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`, http.StatusNotFound},
		{"off grid", `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-24T09:10:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`, http.StatusUnprocessableEntity},
		{"in the past", `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-23T08:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`, http.StatusUnprocessableEntity},
		{"outside window", `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-09-07T09:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`, http.StatusUnprocessableEntity},
		{"bad email", `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-25T09:00:00Z","guestName":"Иван","guestEmail":"nope"}`, http.StatusBadRequest},
		{"bad time", `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-25","guestName":"Иван","guestEmail":"ivan@example.com"}`, http.StatusBadRequest},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			status, body := doJSON(t, http.MethodPost, ts.URL+"/bookings", tc.payload)
			if status != tc.want {
				t.Fatalf("status %d body %s, want %d", status, body, tc.want)
			}
		})
	}
}

func TestCreateBookingAdjacentSlots(t *testing.T) {
	ts := newTestServer(t)
	consult := createEventType(t, ts, consultType)

	mk := func(start string) string {
		return `{"eventTypeId":"` + consult.ID + `","startsAt":"` + start + `","guestName":"Иван","guestEmail":"ivan@example.com"}`
	}
	// 09:00 занято, соседний 09:30 свободен
	if status, body := doJSON(t, http.MethodPost, ts.URL+"/bookings", mk("2026-08-24T09:00:00Z")); status != http.StatusCreated {
		t.Fatalf("first book: status %d body %s", status, body)
	}
	if status, body := doJSON(t, http.MethodPost, ts.URL+"/bookings", mk("2026-08-24T09:30:00Z")); status != http.StatusCreated {
		t.Fatalf("adjacent book: status %d body %s", status, body)
	}
	// повтор того же слота — 409
	if status, body := doJSON(t, http.MethodPost, ts.URL+"/bookings", mk("2026-08-24T09:00:00Z")); status != http.StatusConflict {
		t.Fatalf("same slot again: status %d body %s, want 409", status, body)
	}
}

func TestCreateBookingCrossTypeConflict(t *testing.T) {
	env := newTestEnv(t)
	consult := createEventType(t, env.ts, consultType)
	review := createEventType(t, env.ts, reviewType)

	// бронь на консультацию 10:00-10:30 (26 августа)
	payload := `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-26T10:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`
	if status, body := doJSON(t, http.MethodPost, env.ts.URL+"/bookings", payload); status != http.StatusCreated {
		t.Fatalf("seed: status %d body %s", status, body)
	}

	// слот типа "Разбор кода" 10:00-10:45 пересекает 10:00-10:30 — 409
	payload = `{"eventTypeId":"` + review.ID + `","startsAt":"2026-08-26T10:00:00Z","guestName":"Петя","guestEmail":"petya@example.com"}`
	status, body := doJSON(t, http.MethodPost, env.ts.URL+"/bookings", payload)
	if status != http.StatusConflict {
		t.Fatalf("cross-type overlap: status %d body %s, want 409", status, body)
	}
}

func TestUpcomingAndCancel(t *testing.T) {
	env := newTestEnv(t)
	consult := createEventType(t, env.ts, consultType)

	mk := func(startsAt, name string) string {
		return `{"eventTypeId":"` + consult.ID + `","startsAt":"` + startsAt + `","guestName":"` + name + `","guestEmail":"` + strings.ToLower(name) + `@example.com"}`
	}
	for _, s := range []string{"2026-08-24T10:00:00Z", "2026-08-24T09:00:00Z", "2026-08-25T09:00:00Z"} {
		status, body := doJSON(t, http.MethodPost, env.ts.URL+"/bookings", mk(s, "guest"+s[8:10]))
		if status != http.StatusCreated {
			t.Fatalf("book %s: status %d body %s", s, status, body)
		}
	}
	// прошлое бронирование (вставлено напрямую, чтобы обойти валидацию слота)
	past := domain.Booking{
		ID:          "bk-past",
		EventTypeID: consult.ID,
		StartsAt:    time.Date(2026, 8, 22, 9, 0, 0, 0, time.UTC),
		EndsAt:      time.Date(2026, 8, 22, 9, 30, 0, 0, time.UTC),
		GuestName:   "Прошлый",
		GuestEmail:  "past@example.com",
	}
	if ok, err := env.store.CreateBooking(past); err != nil || !ok {
		t.Fatalf("seed past booking: ok=%v err=%v", ok, err)
	}

	status, body := doJSON(t, http.MethodGet, env.ts.URL+"/admin/bookings", "")
	if status != http.StatusOK {
		t.Fatalf("upcoming: status %d", status)
	}
	bookings := decode[[]BookingDTO](t, body)
	if len(bookings) != 3 {
		t.Fatalf("expected 3 upcoming bookings, got %d: %+v", len(bookings), bookings)
	}
	// сортировка по возрастанию
	for i := 1; i < len(bookings); i++ {
		if bookings[i-1].StartsAt > bookings[i].StartsAt {
			t.Fatalf("not sorted: %+v", bookings)
		}
	}

	// отмена
	target := bookings[0]
	status, _ = doJSON(t, http.MethodDelete, env.ts.URL+"/admin/bookings/"+target.ID, "")
	if status != http.StatusNoContent {
		t.Fatalf("cancel: status %d, want 204", status)
	}
	status, _ = doJSON(t, http.MethodDelete, env.ts.URL+"/admin/bookings/"+target.ID, "")
	if status != http.StatusNotFound {
		t.Fatalf("re-cancel: status %d, want 404", status)
	}
}

func TestDeleteEventTypeCascadesBookings(t *testing.T) {
	ts := newTestServer(t)
	consult := createEventType(t, ts, consultType)

	payload := `{"eventTypeId":"` + consult.ID + `","startsAt":"2026-08-24T09:00:00Z","guestName":"Иван","guestEmail":"ivan@example.com"}`
	if status, _ := doJSON(t, http.MethodPost, ts.URL+"/bookings", payload); status != http.StatusCreated {
		t.Fatal("seed booking failed")
	}

	status, _ := doJSON(t, http.MethodDelete, ts.URL+"/admin/event-types/"+consult.ID, "")
	if status != http.StatusNoContent {
		t.Fatalf("delete type: status %d", status)
	}

	_, body := doJSON(t, http.MethodGet, ts.URL+"/admin/bookings", "")
	if got := decode[[]BookingDTO](t, body); len(got) != 0 {
		t.Fatalf("expected cascade delete, got %d bookings", len(got))
	}
}
