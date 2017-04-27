package sechatapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// writeJson writes the specified JSON value to the client.
func (s *Server) writeJson(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		s.log.Error(err)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// writeError writes an error message to the client.
func (s *Server) writeError(w http.ResponseWriter, message string) {
	s.writeJson(w, map[string]interface{}{
		"error": message,
	})
}
