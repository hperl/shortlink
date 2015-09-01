package main

import (
	"log"
	"net/http"
	"os"
)

func setupMux(db *store) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/admin/delete", adminHandle{store: db})
	mux.Handle("/admin", adminHandle{store: db})
	mux.Handle("/", redirectHandle{store: db})
	return mux
}

func main() {
	storeFile := os.Getenv("STORE_FILE")
	file, err := os.Open(storeFile)
	if err != nil {
		log.Printf("Store File: %v", err)
		file = nil
	} else {
		defer file.Close()
	}
	db := NewStore(file)
	db.filename = storeFile
	mux := setupMux(db)
	panic(http.ListenAndServe(":80", mux))
}
