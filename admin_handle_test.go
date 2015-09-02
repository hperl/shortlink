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
	testStore.DeleteAll()

	type testDatum struct {
		from           string
		to             string
		expectedStatus int
		expectedLength int
	}

	tests := []testDatum{
		{"lgt15", "http://www.google.de", http.StatusOK, 1}, // works first time
		{"lgt15", "http://www.google.de", http.StatusOK, 1}, // does not work again second time
		{"admin", "http://www.google.de", http.StatusOK, 1}, // "admin" is a protected path
		{"fooab", "some.invalid.domains", http.StatusOK, 1}, // invalid domain
		{"datei", "http://www.google.de", http.StatusOK, 1}, // "datei" is a protected path
		{"a/./b", "http://www.google.de", http.StatusOK, 1}, // only alphanumerical
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

func Test_AdminDeleteRedirect(t *testing.T) {
	os.Setenv(ADMIN_USER, "foo")
	os.Setenv(ADMIN_PASSWORD, "bar")

	testStore.DeleteAll()
	if len(testStore.Redirects()) != 0 {
		t.Error("redirect was not cleared")
	}
	testStore.Add(&redirect{"delete-test", "http://www.example.com"})
	if len(testStore.Redirects()) != 1 {
		t.Error("redirect was not added")
	}
	if _, ok := testStore.Get("delete-test"); !ok {
		t.Fatal("redirect was not added")
	}

	req, err := http.NewRequest("GET", server.URL+"/admin/delete?from=delete-test", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.SetBasicAuth("foo", "bar")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if len(testStore.Redirects()) != 0 {
		for _, r := range testStore.Redirects() {
			t.Logf("%+v", r)
		}
		t.Error("redirect was not deleted")
	}
}
