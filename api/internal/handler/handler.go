package handler

import (
	"ama-server/internal/store/pgstore"
	"context"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type ApiHandler struct {
	Q           *pgstore.Queries
	R           *chi.Mux
	Upgrader    websocket.Upgrader
	Subscribers map[string]map[*websocket.Conn]context.CancelFunc
	Mu          *sync.Mutex
}
