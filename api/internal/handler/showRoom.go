package handler

import "net/http"

func (h ApiHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	room, _, _, ok := h.readRoom(w, r)
	if !ok {
		return
	}

	sendJSON(w, room)
}
