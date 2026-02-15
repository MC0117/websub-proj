package subscription

import "sync"

// struct for storing subscriber object
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
	//locks for writing to prevent race condition on s.subscribers map
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers[sub.CallbackURL] = sub
}

func (s *Store) GetSubscribers() []Subscriber {
	//read locks to prevent race condition on s.subscribers map
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Subscriber

	for _, sub := range s.subscribers {
		result = append(result, sub)
	}
	return result
}

func (s *Store) GetSubscribersByTopic(topic string) []Subscriber {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []Subscriber

	for _, sub := range s.subscribers {
		if sub.Topic == topic {
			result = append(result, sub)
		}
	}
	return result
}
