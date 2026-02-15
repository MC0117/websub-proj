package handlers

import (
	"fmt"
	"hub/internal/subscription"
	"math/rand"
	"net/http"
	"strconv"
)

type SubscriptionHandler struct {
	store *subscription.Store
}

// "constructor" that is called externally from main.go
func NewSubscriptionHandler(store *subscription.Store) *SubscriptionHandler {
	return &SubscriptionHandler{store: store}
}

// method that "serves" the endpoint
func (h *SubscriptionHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "HTTP method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(rw, "Bad Request", http.StatusBadRequest)
		return
	}

	callback := r.FormValue("hub.callback")
	secret := r.FormValue("hub.secret")
	topic := r.FormValue("hub.topic")
	mode := r.FormValue("hub.mode")
	// todo: should probably validate these arent empty, return 400 if they are

	challenge := strconv.Itoa(rand.Intn(1000000))

	verifyURL := fmt.Sprintf("%s?hub.mode=%s&hub.topic=%s&hub.challenge=%s", callback, mode, topic, challenge)

	resp, err := http.Get(verifyURL)
	// todo: need to check err before accessing resp, will panic if http.Get fails
	// also move the defer resp.Body.Close() after the err check

	fmt.Printf("Verify URL: %s\n", verifyURL)
	fmt.Printf("Response Status: %d\n", resp.StatusCode)

	defer resp.Body.Close()

	if err != nil {
		http.Error(rw, "Validation of Subscription Failed", http.StatusNotFound)
		return
	}

	if resp.StatusCode == http.StatusOK {
		newSub := subscription.Subscriber{
			CallbackURL: callback,
			Secret:      secret,
			Topic:       topic,
		}
		h.store.Add(newSub)
		rw.WriteHeader(http.StatusAccepted)
		fmt.Printf("Added subscriber: %s\n", callback)
	}
}
