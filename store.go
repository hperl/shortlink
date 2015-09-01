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
	redirects []*redirect
	filename  string
}

func (s *store) Redirects() []*redirect {
	return s.redirects
}

func (s *store) Add(redirect *redirect) error {
	for _, r := range s.redirects {
		if r.From == redirect.From {
			return fmt.Errorf("Quelle %q existiert schon", redirect.From)
		}
	}
	s.redirects = append(s.redirects, redirect)
	s.writeFile()

	return nil
}

func (s *store) writeFile() {
	data, _ := json.Marshal(s.redirects)
	if err := ioutil.WriteFile(s.filename, data, 0644); err != nil {
		log.Printf("Error writing data to %q: %v.", s.filename, err)
	}
}

func (s *store) DeleteAll() {
	s.redirects = make([]*redirect, 0)
}

func NewStore(data io.Reader) (s *store) {
	s = &store{filename: os.Getenv("STORE_FILE")}
	bytes, err := ioutil.ReadAll(data)
	if err != nil {
		log.Print(err)
		return
	}
	if err = json.Unmarshal(bytes, &s.redirects); err != nil {
		log.Print(err)
		return
	}

	return
}
