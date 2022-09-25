package ping

import (
	"encoding/json"
	"net/http"
)

// Get método para 'health check'
func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode("pong")
}
