package conversation

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

var ErrNotFound = errors.New("conversation not found")

type Filter struct {
	Status   Status
	Priority Priority
	Search   string
}

func (f Filter) matches(c Conversation) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Priority != "" && c.Priority != f.Priority {
		return false
	}
	if f.Search != "" {
		q := strings.ToLower(f.Search)
		match := strings.Contains(strings.ToLower(c.CustomerName), q) ||
			strings.Contains(strings.ToLower(c.CustomerEmail), q) ||
			strings.Contains(strings.ToLower(c.Subject), q)
		if !match {
			return false
		}
	}
	return true
}

type Update struct {
	Status   *Status
	Priority *Priority
}

type Store struct {
	mu            sync.RWMutex
	conversations []Conversation
}

func NewStore(seed []Conversation) *Store {
	conversations := make([]Conversation, len(seed))
	copy(conversations, seed)
	return &Store{conversations: conversations}
}

func (s *Store) List(f Filter) []Conversation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Conversation, 0, len(s.conversations))
	for _, c := range s.conversations {
		if f.matches(c) {
			result = append(result, c)
		}
	}
	return result
}

func (s *Store) Get(id string) (Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.conversations {
		if c.ID == id {
			return c, nil
		}
	}
	return Conversation{}, fmt.Errorf("%w: %s", ErrNotFound, id)
}

func (s *Store) Update(id string, u Update) (Conversation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.conversations {
		if s.conversations[i].ID != id {
			continue
		}
		if u.Status != nil {
			s.conversations[i].Status = *u.Status
		}
		if u.Priority != nil {
			s.conversations[i].Priority = *u.Priority
		}
		return s.conversations[i], nil
	}
	return Conversation{}, fmt.Errorf("%w: %s", ErrNotFound, id)
}
