package main

import (
	"log"
	"net/http"
)

type redirectHandle struct {
	store *store
}

func (h redirectHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, r := range h.store.Redirects() {
		if req.URL.String() == "/"+r.From {
			http.Redirect(w, req, r.To, http.StatusTemporaryRedirect)
			return
		}
	}
	log.Printf("%v: redirect not found", req.URL)
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Der Link konnte nicht gefunden werden."))
}
