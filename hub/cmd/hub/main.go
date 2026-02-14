package main

import (
	"hub/internal/handlers"
	"hub/internal/subscription"
	"log"
	"net/http"
)

func main() {

	store := subscription.NewStore()

	// init handlers
	subHandler := handlers.NewSubscriptionHandler(store)
	pubHandler := handlers.NewPublishHandler(store)

	// init Routes
	http.Handle("/subscribe", subHandler)
	http.Handle("/publish", pubHandler)

	log.Println("Hub listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
