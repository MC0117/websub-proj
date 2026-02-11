package handlers

import (
	"hub/internal/delivery"
	"hub/internal/subscription"
	"io"
	"net/http"
)

type PublishHandler struct {
	store *subscription.Store
}

func NewPublishHandler(s *subscription.Store) *PublishHandler {
	return &PublishHandler{}
}

func (h *PublishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "invalid message body", http.StatusBadRequest)
		return
	}

	subscribers := h.store.GetSubscribers()

	for _, sub := range subscribers {
		//todo optimizable? maybe separate into two sections
		go delivery.SendPayload(sub.CallbackURL, sub.Secret, body)
	}

}
