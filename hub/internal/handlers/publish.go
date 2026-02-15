package handlers

import (
	"hub/internal/delivery"
	"hub/internal/subscription"
	"io"
	"log"
	"net/http"
)

type PublishHandler struct {
	store *subscription.Store
}

func NewPublishHandler(s *subscription.Store) *PublishHandler {
	return &PublishHandler{store: s}
}

func (h *PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	log.Printf("Message Body: %s", string(body))
	if err != nil || len(body) == 0 {
		http.Error(w, "invalid message body", http.StatusBadRequest)
		return
	}

	topic := r.URL.Query().Get("hub.topic")
	if topic == "" {
		http.Error(w, "missing hub.topic query parameter", http.StatusBadRequest)
		return
	}

	subscribers := h.store.GetSubscribersByTopic(topic)

	for _, sub := range subscribers {
		go func(s subscription.Subscriber) {
			if err := delivery.SendPayload(s.CallbackURL, s.Secret, body); err != nil {
				log.Printf("Delivery Failed for %s: %v\n", s.CallbackURL, err)
			}
		}(sub)
	}
	log.Printf("Published to %d subscriber(s) on topic %s", len(subscribers), topic)
}
