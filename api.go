package sechatapi

import (
	"encoding/json"
	"net/http"
)

type sendParams struct {
	Room int    `json:"room"`
	Text string `json:"text"`
}

func (s *Server) handleSend(w http.ResponseWriter, r *http.Request) {
	v := &sendParams{}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		s.writeError(w, err.Error())
		return
	}
	if err := s.conn.Send(v.Room, v.Text); err != nil {
		s.writeError(w, err.Error())
		return
	}
	s.writeJson(w, map[string]interface{}{})
}
