package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type store struct {
	redirects map[string]*redirect
	filename  string
}

func (s *store) Redirects() (redirects []*redirect) {
	if s.redirects == nil {
		return
	}

	for _, r := range s.redirects {
		redirects = append(redirects, r)
	}

	return
}

func (s *store) Add(redirect *redirect) error {
	if _, ok := s.redirects[redirect.From]; ok {
		return fmt.Errorf("Quelle %q existiert schon", redirect.From)
	} else {
		if s.redirects == nil {
			s.DeleteAll()
		}
		s.redirects[redirect.From] = redirect
	}

	s.writeFile()

	return nil
}

func (s *store) Get(from string) (*redirect, bool) {
	r, ok := s.redirects[from]
	return r, ok
}

func (s *store) writeFile() {
	data, _ := json.Marshal(s.Redirects())
	if err := ioutil.WriteFile(s.filename, data, 0644); err != nil {
		log.Printf("Error writing data to %q: %v.", s.filename, err)
	}
}

func (s *store) Delete(from string) {
	delete(s.redirects, from)
}

func (s *store) DeleteAll() {
	s.redirects = make(map[string]*redirect)
}

func NewStore(data io.Reader) (s *store) {
	s = &store{
		filename:  os.Getenv("STORE_FILE"),
		redirects: make(map[string]*redirect),
	}

	var redirects []*redirect

	if data != nil {
		bytes, err := ioutil.ReadAll(data)
		if err != nil {
			log.Print(err)
			return
		}
		if err = json.Unmarshal(bytes, &redirects); err != nil {
			log.Print(err)
			return
		}
		for _, r := range redirects {
			s.redirects[r.From] = r
		}
	}

	return
}
