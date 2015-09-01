package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_GetAdminFailsWithoutPassword(t *testing.T) {
	os.Setenv(ADMIN_USER, "")
	os.Setenv(ADMIN_PASSWORD, "")
	req, err := http.NewRequest("GET", server.URL+"/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("", "")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected status %v, got %v", http.StatusUnauthorized, res.StatusCode)
	}
}

func Test_GetAdminSucceedsWithPassword(t *testing.T) {
	os.Setenv(ADMIN_USER, "foo")
	os.Setenv(ADMIN_PASSWORD, "bar")
	req, err := http.NewRequest("GET", server.URL+"/admin", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("foo", "bar")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, res.StatusCode)
	}
}

func Test_AdminPostNewRedirect(t *testing.T) {
	type testDatum struct {
		from           string
		to             string
		expectedStatus int
		expectedLength int
	}

	tests := []testDatum{
		{"lgt15", "http://www.google.de", http.StatusOK, 1},       // works first time
		{"lgt15", "http://www.google.de", http.StatusNotFound, 1}, // does not work again second time
		{"admin", "http://www.google.de", http.StatusNotFound, 1}, // "admin" is a protected path
		{"foo", "invalid.domain", http.StatusNotFound, 1},         // "admin" is a protected path
	}

	os.Setenv(ADMIN_USER, "foo")
	os.Setenv(ADMIN_PASSWORD, "bar")

	for _, test := range tests {
		data := url.Values{
			"from": {test.from},
			"to":   {test.to},
		}
		req, err := http.NewRequest("POST", server.URL+"/admin", strings.NewReader(data.Encode()))
		if err != nil {
			t.Fatal(err)
		}
		req.SetBasicAuth("foo", "bar")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		if res.StatusCode != test.expectedStatus {
			t.Errorf("expected status %v, got %v", test.expectedStatus, res.StatusCode)
		}
		if len(testStore.Redirects()) != test.expectedLength {
			t.Errorf("Expected exactly %d redirect(s), was %d", test.expectedLength, len(testStore.Redirects()))
		}
	}
}
