package subscription

import "sync"

// storing this
type Subscriber struct {
	CallbackURL string
	Secret      string
	Topic       string
}

type Store struct {
	mu          sync.RWMutex
	subscribers map[string]Subscriber
}

func NewStore() *Store {
	return &Store{
		subscribers: make(map[string]Subscriber),
	}
}

func (s *Store) Add(sub Subscriber) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[sub.CallbackURL] = sub
}

func (s *Store) GetSubscribers() []Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Subscriber

	for i := 0; i < len(s.subscribers); i++ {

	}

}
