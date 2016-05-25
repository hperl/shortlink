package main

import (
	"encoding/json"
	"log"

	"gopkg.in/redis.v3"
)

const (
	redisRedirectsKey = "shortlink_redirects"
)

type store struct {
	db *redis.Client
}

func (s *store) Redirects() (redirects []*redirect) {
	results, err := s.db.HGetAllMap(redisRedirectsKey).Result()
	if err != nil {
		log.Printf("store.Redirects(): %v", err)
	}

	for _, jsonString := range results {
		redirect := new(redirect)
		json.Unmarshal([]byte(jsonString), redirect)
		redirects = append(redirects, redirect)
	}

	return
}

func (s *store) Add(redirect *redirect) error {
	jsonBytes, err := json.Marshal(redirect)
	if err != nil {
		panic(err)
	}
	s.db.HSet(redisRedirectsKey, redirect.From, string(jsonBytes)).Result()

	return nil
}

func (s *store) Get(from string) (r *redirect, ok bool) {
	jsonString, err := s.db.HGet(redisRedirectsKey, from).Result()
	if err != nil {
		log.Printf("store.Get(%q): %v", from, err)
		return nil, false
	}
	r = new(redirect)
	if err := json.Unmarshal([]byte(jsonString), r); err != nil {
		log.Printf("store.Get(%q): %v", from, err)
		return nil, false
	}

	return r, true
}

func (s *store) Delete(from string) {
	cnt, err := s.db.HDel(redisRedirectsKey, from).Result()
	if err != nil {
		log.Printf("store.Delete(%q): %v", from, err)
	}
	if cnt != 1 {
		log.Printf("store.Delete(%q): not found", from)
	}
}

func (s *store) DeleteAll() {
	s.db.Del(redisRedirectsKey)
}

func (s *store) Close() error {
	return s.db.Close()
}

/// NewStore creates a new redirect store.
func NewStore() (s *store) {
	s = &store{
		db: redis.NewClient(&redis.Options{Addr: "redis:6379"}),
	}

	return
}
