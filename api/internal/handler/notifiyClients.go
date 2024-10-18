package handler

import "log/slog"

type Message struct {
	Kind   string `json:"kind"`
	Value  any    `json:"value"`
	RoomID string `json:"-"`
}

func (h ApiHandler) NotifyClients(msg Message) {
	h.Mu.Lock()
	defer h.Mu.Unlock()

	subscribers, ok := h.Subscribers[msg.RoomID]
	if !ok || len(subscribers) == 0 {
		return
	}
	for coon, cancel := range subscribers {
		if err := coon.WriteJSON(msg); err != nil {
			slog.Error("failed to send message to client", "error", err)
			cancel()
		}
	}
}
