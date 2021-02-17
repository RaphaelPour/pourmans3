package main

import (
	"errors"
	"fmt"
	"hash/maphash"
)

type Storage struct {
	hash    maphash.Hash
	kvStore map[string]string
}

func NewStorage() Storage {
	return Storage{
		kvStore: map[string]string{},
	}
}

func (s *Storage) Get(key string) (string, error) {
	if key == "" {
		return "", errors.New("empty key")
	}

	if value, ok := s.kvStore[key]; ok {
		return value, nil
	}

	return "", fmt.Errorf("unknown key '%s'", key)
}

func (s *Storage) Set(value string) string {
	s.hash.Reset()
	// Docu: WriteString never fails
	_, _ = s.hash.WriteString(value)
	key := fmt.Sprintf("%x", s.hash.Sum64())
	s.kvStore[key] = value
	return key
}

func (s *Storage) Delete(key string) {
	delete(s.kvStore, key)
}

func (s *Storage) All() map[string]string {
	return s.kvStore
}
