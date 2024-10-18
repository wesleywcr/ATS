package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5"
)

func (h ApiHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	rawRoomID := chi.URLParam(r, "room_id")
	roomID, err := uuid.Parse(rawRoomID)
	if err != nil {
		http.Error(w, "invalid room id", http.StatusBadRequest)
		return
	}

	_, err = h.Q.GetRoom(r.Context(), roomID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "room not found", http.StatusBadRequest)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	c, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Warn("failed to upgrade connection", "error", err)
		http.Error(w, "failed to upgrade ws connection", http.StatusBadRequest)
		return
	}
	defer c.Close()
	ctx, cancel := context.WithCancel(r.Context())
	h.Mu.Lock()
	if _, ok := h.Subscribers[rawRoomID]; !ok {
		h.Subscribers[rawRoomID] = make(map[*websocket.Conn]context.CancelFunc)
	}
	slog.Info("new client connected", "room_id", rawRoomID, "client_ip", r.RemoteAddr)
	h.Subscribers[rawRoomID][c] = cancel
	h.Mu.Unlock()

	<-ctx.Done()

	h.Mu.Lock()
	delete(h.Subscribers[rawRoomID], c)
	h.Mu.Unlock()

}
