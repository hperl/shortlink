package main

import (
	"strings"
	"testing"
)

func Test_NewStoreIsReady(t *testing.T) {
	store := NewStore(nil)
	if len(store.Redirects()) != 0 {
		t.Error("Empty store should have no redirects")
	}
	store.Add(&redirect{"a", "b"})
	if len(store.Redirects()) != 1 {
		t.Error("Empty store should have 1 redirect")
	}

	store.Delete("a")
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
