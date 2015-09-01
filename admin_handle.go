package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

import "os"

const ADMIN_USER = "ADMIN_USER"
const ADMIN_PASSWORD = "ADMIN_PASSWORD"

type adminHandle struct {
	store  *store
	errMsg string
}

func (h adminHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user := os.Getenv(ADMIN_USER)
	password := os.Getenv(ADMIN_PASSWORD)

	u, p, ok := req.BasicAuth()
	if !ok || user == "" || password == "" || u != user || p != password {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"Link-Administration\"")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Benutzername und Passwort bekommst du von Henning."))
		return
	}

	if req.Method == "POST" {
		if err := h.addRedirect(w, req); err != nil {
			h.errMsg = err.Error()
		}
	}
	h.renderAdminInterface(w)
}

func (h *adminHandle) addRedirect(w http.ResponseWriter, req *http.Request) (err error) {
	from := req.PostFormValue("from")
	to := req.PostFormValue("to")

	if err := h.validateTo(to); err != nil {
		return err
	}
	if err := h.validateFrom(from); err != nil {
		return err
	}

	if err := h.store.Add(&redirect{From: from, To: to}); err != nil {
		return err
	}

	return nil
}

func (h *adminHandle) validateTo(to string) error {
	if to == "" {
		return errors.New("Ziel ist leer.")
	}
	if !isUrlReachable(to) {
		return fmt.Errorf("%q ist keine gültige URL.", to)
	}
	return nil
}

func isUrlReachable(url string) bool {
	_, err := http.Get(url)
	return err == nil
}

func (h *adminHandle) validateFrom(from string) error {
	if from == "" {
		return errors.New("Kurzlink ist leer.")
	}
	if from == "admin" {
		return errors.New("'admin' darf nicht als Kurzlink verwendet werden.")
	}
	return nil
}

func (h *adminHandle) renderAdminInterface(w http.ResponseWriter) {
	t, err := template.New("index").Parse(adminIndexTemplate)
	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		Title     string
		BaseURL   string
		Redirects []*redirect
		ErrorMsg  string
	}{
		Title:     "Link-Administration",
		BaseURL:   os.Getenv("BASE_URL"),
		Redirects: h.store.Redirects(),
		ErrorMsg:  h.errMsg,
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
		return
	}
}

const adminIndexTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		<h1>{{.Title}}</h1>
		<h2>Aktuelle Kurzlinks</h2>
		<ul>
			{{range .Redirects}}
				<li>
					<a href="{{$.BaseURL}}/{{.From}}">{{$.BaseURL}}/{{.From}}</a>
					→
					<a href="{{.To}}">{{.To}}</a>
				</li>
			{{else}}
				<li><i>keine Redirects</i></li>
			{{end}}
		</ul>
		<hr>
		<h2>Neuer Kurzlink</h2>
		{{if .ErrorMsg}}
			<p><strong>{{.ErrorMsg}}</strong></p>
		{{end}}
		<form action="/admin" method="post">
			Kurzlink:<br>
			<input type="text" name="from">
			<br>
			Zieladresse:<br>
			<input type="text" name="to">
			<br>
			<input type="submit" value="Erstellen">
		</form>
	</body>
</html>`
