package handlers

import (
	"encoding/json"
	"net/http"
)

func encodeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func decodeJSON(r *http.Request, target any) error {
	return json.NewDecoder(r.Body).Decode(target)
}

func writeError(w http.ResponseWriter, status int, message string) {
	encodeJSON(w, status, map[string]string{"error": message})
}
