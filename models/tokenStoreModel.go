package models

import (
	"sync"
)

type SafeTokenStore struct {
	store map[uint]string
	mutex sync.RWMutex
}

func NewSafeTokenStore() *SafeTokenStore {
	return &SafeTokenStore{
		store: make(map[uint]string),
	}
}

func (s *SafeTokenStore) Set(userID uint, token string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.store[userID] = token
}

func (s *SafeTokenStore) Get(userID uint) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	token, exists := s.store[userID]
	return token, exists
}

func (s *SafeTokenStore) Delete(userID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.store, userID)
}

func (s *SafeTokenStore) Validate(userID uint, token string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	storedToken, exists := s.store[userID]
	return exists && storedToken == token
}
