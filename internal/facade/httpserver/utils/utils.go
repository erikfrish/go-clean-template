package utils

import "net/http"

func Write200(body []byte, w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(body)
}

func Write400(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}

func Write500(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}

func Write204(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func Write401(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}

func Write403(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusForbidden)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(err))
}
