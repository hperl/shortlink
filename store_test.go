package main

import (
	"strings"
	"testing"
)

func Test_NewStoreIsEmpty(t *testing.T) {
	store := new(store)
	if len(store.Redirects()) != 0 {
		t.Error("Empty store should have no redirects")
	}
}

func Test_ReadStoreState(t *testing.T) {
	s := strings.NewReader(`[{"From":"foo","To":"http://yfu.de"}]`)

	store := NewStore(s)
	if len(store.Redirects()) != 1 {
		t.Fatal("Store should have one redirect")
	}
	r := store.Redirects()[0]

	if r.From != "foo" || r.To != "http://yfu.de" {
		t.Error(r)
	}
}
