package api

import (
	"ama-server/internal/handler"
	"ama-server/internal/store/pgstore"
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
)

func NewHandler(q *pgstore.Queries) http.Handler {
	a := handler.ApiHandler{
		Q: q,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Subscribers: make(map[string]map[*websocket.Conn]context.CancelFunc),
		Mu:          &sync.Mutex{},
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Recoverer, middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/subscribe/{room_id}", a.Subscribe)

	r.Route("/api", func(r chi.Router) {
		r.Route("/rooms", func(r chi.Router) {
			r.Post("/", a.CreateRoom)
			r.Get("/", a.GetRooms)

			r.Route("/{room_id}", func(r chi.Router) {
				r.Get("/", a.GetRoom)

				r.Route("/messages", func(r chi.Router) {
					r.Post("/", a.CreateRoomMessage)
					r.Get("/", a.GetRoomMessages)

					r.Route("/{message_id}", func(r chi.Router) {
						r.Get("/", a.GetRoomMessage)
						r.Patch("/react", a.ReactToMessage)
						r.Delete("/react", a.RemoveReactFromMessage)
						r.Patch("/answer", a.MarkMessageAsAnswer)

					})
				})
			})
		})
	})

	a.R = r
	return a
}
