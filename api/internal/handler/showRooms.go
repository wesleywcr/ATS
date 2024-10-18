package handler

import (
	"ama-server/internal/store/pgstore"
	"log/slog"
	"net/http"
)

func (h ApiHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.Q.GetRooms(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		slog.Error("failed to get rooms", "error", err)
		return
	}

	if rooms == nil {
		rooms = []pgstore.Room{}
	}

	sendJSON(w, rooms)
}
