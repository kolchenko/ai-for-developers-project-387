package api

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"callcalendar/backend/internal/store"
)

type Server struct {
	store  *store.Store
	now    func() time.Time
	webDir string
}

func NewServer(st *store.Store) http.Handler {
	return newServer(st, time.Now, "")
}

func NewServerWithWeb(st *store.Store, webDir string) http.Handler {
	return newServer(st, time.Now, webDir)
}

func newServer(st *store.Store, now func() time.Time, webDir string) http.Handler {
	return (&Server{store: st, now: now, webDir: webDir}).routes()
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	// Каждый маршрут регистрируется дважды: без префикса (для dev и тестов,
	// когда vite-прокси срезает `/api`) и с префиксом `/api` (для продакшена,
	// где SPA и API живут за одним хостом).
	api := func(method, path string, h http.HandlerFunc) {
		mux.HandleFunc(method+" "+path, h)
		mux.HandleFunc(method+" /api"+path, h)
	}

	api("GET", "/event-types", s.listEventTypes)
	api("GET", "/event-types/{eventTypeId}/slots", s.getSlots)
	api("POST", "/bookings", s.createBooking)
	api("POST", "/admin/login", s.adminLogin)
	api("POST", "/admin/event-types", s.adminCreateEventType)
	api("PATCH", "/admin/event-types/{eventTypeId}", s.adminUpdateEventType)
	api("DELETE", "/admin/event-types/{eventTypeId}", s.adminDeleteEventType)
	api("GET", "/admin/bookings", s.adminUpcomingBookings)
	api("DELETE", "/admin/bookings/{bookingId}", s.adminCancelBooking)

	if s.webDir != "" {
		mux.Handle("GET /", s.spaHandler())
	}

	return withLogging(mux)
}

func NewTestServer(st *store.Store, now func() time.Time) http.Handler {
	return newServer(st, now, "")
}

// spaHandler отдаёт статические файлы собранного фронтенда, а для всех
// остальных путей возвращает index.html (SPA-fallback).
func (s *Server) spaHandler() http.Handler {
	fs := http.FileServer(http.Dir(s.webDir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "..") {
			http.NotFound(w, r)
			return
		}
		path := filepath.Join(s.webDir, strings.TrimPrefix(r.URL.Path, "/"))
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, filepath.Join(s.webDir, "index.html"))
	})
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}
