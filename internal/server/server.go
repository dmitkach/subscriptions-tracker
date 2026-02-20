package server

import (
	"log"
	"net/http"
	"subscriptions-tracker/internal/handlers"

	"github.com/gorilla/mux"
)

type Server struct {
	Server http.Server
}

func RunServer(h *handlers.Handler) {
	r := mux.NewRouter()

	r.HandleFunc("/subscriptions", h.CreateSubscription).Methods(http.MethodPost)
	r.HandleFunc("/subscriptions/{id}", h.GetSubscription).Methods(http.MethodGet)
	r.HandleFunc("/users/{userID}/subscriptions", h.GetUserSubscriptions).Methods(http.MethodGet)

	if err := http.ListenAndServe("0.0.0.0:8080", r); err != nil {
		log.Fatal(err)
	}
	log.Println("Server is running")
}
