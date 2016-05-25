package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func Test_GetAdminFailsWithoutPassword(t *testing.T) {
	os.Setenv(adminUserEnv, "")
	os.Setenv(adminPassEnv, "")
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
	os.Setenv(adminUserEnv, "foo")
	os.Setenv(adminPassEnv, "bar")
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

	type testData struct {
		from              string
		to                string
		expectedRedirects int
	}

	tests := []testData{
		{"lgt15", "http://www.google.de", 1},  // works first time
		{"lgt-15", "http://www.google.de", 1}, // works with dashes
		{"admin", "http://www.google.de", 0},  // "admin" is a protected path
		{"fooab", "some.invalid.domains", 0},  // invalid domain
		{"datei", "http://www.google.de", 0},  // "datei" is a protected path
		{"a/./b", "http://www.google.de", 0},  // only alphanumerical
	}

	os.Setenv(adminUserEnv, "foo")
	os.Setenv(adminPassEnv, "bar")

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
		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, res.StatusCode)
		}
		if len(testStore.Redirects()) != test.expectedRedirects {
			t.Logf("Test data %+v", test)
			t.Errorf("Expected exactly %d redirect(s), was %d",
				test.expectedRedirects,
				len(testStore.Redirects()),
			)
		}

		testStore.DeleteAll()
	}
}

func Test_validateFrom(t *testing.T) {
	type testData struct {
		string      string
		expectValid bool
	}

	tests := []testData{
		{"lgt15", true},
		{"lgt-15", true},
		{"-+_", true},
		{"../invalid", false},
	}

	for _, test := range tests {
		err := validateFrom(test.string)
		t.Logf("test data %+v", test)
		if test.expectValid && err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !test.expectValid && err == nil {
			t.Error("expected error, got none")
		}
	}
}

func Test_AdminDeleteRedirect(t *testing.T) {
	os.Setenv(adminUserEnv, "foo")
	os.Setenv(adminPassEnv, "bar")

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
