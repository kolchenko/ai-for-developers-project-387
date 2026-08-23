package api

import (
	"encoding/json"
	"net/http"
)

type errorBody struct {
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorBody{Message: message})
}

func writeNotFound(w http.ResponseWriter) {
	writeError(w, http.StatusNotFound, "Ресурс не найден")
}

func writeBadRequest(w http.ResponseWriter, message string) {
	writeError(w, http.StatusBadRequest, message)
}

func writeInternal(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "Внутренняя ошибка сервера")
}
