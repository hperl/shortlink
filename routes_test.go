package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
)
import "testing"

var server *httptest.Server
var testStore *store

func TestMain(m *testing.M) {
	setupStore()
	setupServer()
	retVal := m.Run()
	server.Close()
	os.Remove(testStore.filename)
	os.Exit(retVal)
}

func setupStore() {
	testStore = new(store)
	f, err := ioutil.TempFile("", "")
	if err != nil {
		panic(err)
	}
	testStore.filename = f.Name()
	f.Close()
}

func setupServer() {
	mux := http.NewServeMux()
	mux.Handle("/admin", adminHandle{store: testStore})
	mux.Handle("/", redirectHandle{store: testStore})
	server = httptest.NewServer(mux)
}
