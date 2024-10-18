package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (h ApiHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	type _body struct {
		Theme string `json:"theme"`
	}

	var body _body
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	roomID, err := h.Q.InsertRoom(r.Context(), body.Theme)
	if err != nil {
		slog.Error("failed to insert room", "error", err)
		http.Error(w, "somenthing went wrong", http.StatusInternalServerError)
		return
	}

	type response struct {
		ID string `json:"id"`
	}
	data, _ := json.Marshal(response{ID: string(roomID.String())})
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(data)
}
