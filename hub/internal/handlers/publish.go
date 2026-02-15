package handlers

import (
	"fmt"
	"hub/internal/delivery"
	"hub/internal/subscription"
	"io"
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
	fmt.Printf("Message Body: %s\n", string(body))
	if err != nil || len(body) == 0 {
		http.Error(w, "invalid message body", http.StatusBadRequest)
		return
	}

	subscribers := h.store.GetSubscribers()

	for _, sub := range subscribers {
		//todo optimizable? maybe separate into two sections
		// todo: this sends to everyone, should filter by topic
		go delivery.SendPayload(sub.CallbackURL, sub.Secret, body)
	}
	fmt.Printf("Message Sent to all Subscribers: %s\n", string(body))
}
