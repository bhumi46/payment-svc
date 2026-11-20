// Package items implements a small items API with an in-memory store.
package items

import (
	"errors"
	"sort"
	"strconv"
	"sync"
)

// Item is a named thing.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ErrNotFound is returned when an item does not exist.
var ErrNotFound = errors.New("item not found")

// Store is a concurrency-safe in-memory item store.
type Store struct {
	mu     sync.RWMutex
	items  map[string]Item
	nextID int
}

// NewStore returns an empty store.
func NewStore() *Store {
	return &Store{items: map[string]Item{}, nextID: 1}
}

// List returns all items sorted by numeric ID.
func (s *Store) List() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Item, 0, len(s.items))
	for _, it := range s.items {
		out = append(out, it)
	}
	sort.Slice(out, func(i, j int) bool {
		a, _ := strconv.Atoi(out[i].ID)
		b, _ := strconv.Atoi(out[j].ID)
		return a < b
	})
	return out
}

// Get returns one item or ErrNotFound.
func (s *Store) Get(id string) (Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	it, ok := s.items[id]
	if !ok {
		return Item{}, ErrNotFound
	}
	return it, nil
}

// Create assigns an ID and stores the item.
func (s *Store) Create(it Item) Item {
	s.mu.Lock()
	defer s.mu.Unlock()
	it.ID = strconv.Itoa(s.nextID)
	s.nextID++
	s.items[it.ID] = it
	return it
}
