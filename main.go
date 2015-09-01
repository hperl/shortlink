package main

import (
	"net/http"
	"os"
)

var db *store

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/admin", adminHandle{store: db})
	mux.Handle("/", redirectHandle{store: db})
	return mux
}

func main() {
	db = &store{filename: os.Getenv("STORE_FILE")}
	mux := setupMux()
	panic(http.ListenAndServe(":80", mux))
}
