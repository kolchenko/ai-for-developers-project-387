package api

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"net/http"

	"callcalendar/backend/internal/store"
)

func (s *Server) adminLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if !decodeBody(w, r, &req) {
		return
	}
	if req.Username == "" || req.Password == "" {
		writeBadRequest(w, "username и password обязательны")
		return
	}

	admin, err := s.store.GetAdmin(req.Username)
	if err != nil {
		if err == store.ErrNotFound {
			writeError(w, http.StatusUnauthorized, "Неверный логин или пароль")
			return
		}
		writeInternal(w)
		return
	}

	sum := md5.Sum([]byte(req.Password))
	if subtle.ConstantTimeCompare([]byte(hex.EncodeToString(sum[:])), []byte(admin.PasswordHash)) != 1 {
		writeError(w, http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	writeJSON(w, http.StatusOK, LoginResponse{Username: admin.Username})
}
