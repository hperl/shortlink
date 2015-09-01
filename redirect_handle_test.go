package main

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func Test_RedirectHandleRedirects(t *testing.T) {
	to := "http://www.yfu.de"
	testStore.Add(&redirect{From: "foo", To: to})

	client := new(http.Client)
	var redirectURL *url.URL
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		redirectURL = req.URL
		return errors.New("") // we don't want to carry out the redirect, just get the target URL
	}
	client.Get(server.URL + "/foo")

	if redirectURL == nil {
		t.Fatalf("no redirect observed!")
	}
	if redirectURL.String() != to {
		t.Errorf("expected %q, got %q", to, redirectURL.String())
	}
}

func Test_RedirectHandleRedirectsNotFound(t *testing.T) {
	res, err := http.Get(server.URL + "/invalid")
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("expected NotFound, got %v", res.Status)
	}
}
