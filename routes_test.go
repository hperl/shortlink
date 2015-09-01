package main

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
)
import "testing"

var server *httptest.Server
var testStore *store

func TestMain(m *testing.M) {
	testStore = setupStore()
	setupServer(testStore)
	retVal := m.Run()
	server.Close()
	os.Remove(testStore.filename)
	os.Exit(retVal)
}

func setupStore() *store {
	db := NewStore(nil)

	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	db.filename = f.Name()
	f.Close()

	return db
}

func setupServer(db *store) {
	mux := setupMux(db)
	server = httptest.NewServer(mux)
}
