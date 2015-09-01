package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

import "os"

const ADMIN_USER = "ADMIN_USER"
const ADMIN_PASSWORD = "ADMIN_PASSWORD"

type adminHandle struct {
	store   *store
	message string
}

func (h adminHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	user := os.Getenv(ADMIN_USER)
	password := os.Getenv(ADMIN_PASSWORD)

	u, p, ok := req.BasicAuth()
	if !ok || user == "" || password == "" || u != user || p != password {
		w.Header().Set("WWW-Authenticate", `Basic realm="Link-Administration"`)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Benutzername und Passwort bekommst du von Henning."))
		return
	}

	if req.Method == "POST" {
		if err := h.addRedirect(w, req); err != nil {
			h.message = err.Error()
		}
	}
	if req.Method == "GET" && h.isDeleteAction(req.URL) {
		h.message = h.deleteRedirect(req.URL)
		http.Redirect(w, req, "/admin", http.StatusTemporaryRedirect)
	}
	h.renderAdminInterface(w)
}

func (h *adminHandle) isDeleteAction(url *url.URL) bool {
	return strings.Contains(url.Path, "/admin/delete")
}

func (h *adminHandle) deleteRedirect(url *url.URL) string {
	from := url.Query().Get("from")
	if r, ok := h.store.Get(from); ok {
		h.store.Delete(from)
		return fmt.Sprintf("%q → %q wurde gelöscht.", r.From, r.To)
	}
	return fmt.Sprintf("%q konnte nicht gelöscht werden", from)
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
		Message   string
	}{
		Title:     "Link-Administration",
		BaseURL:   os.Getenv("BASE_URL"),
		Redirects: h.store.Redirects(),
		Message:   h.message,
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
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
		<title>{{.Title}}</title>
	</head>
	<body>
		<div class="container">
			<div class="row">
				<div class="col-md-12">
					<h1>{{.Title}}</h1>
					{{if .Message}}
						<p class="bg-primary">
							{{.Message}}
						</p>
					{{end}}
				</div>

				<div class="col-md-6">
					<h2>Aktuelle Kurzlinks</h2>
					<ul class="list-unstyled">
						{{range .Redirects}}
							<li>
								<a href="{{$.BaseURL}}/{{.From}}">{{$.BaseURL}}/{{.From}}</a>
								→
								<a href="{{.To}}">{{.To}}</a>
								<a class="btn btn-default btn-xs" href="/admin/delete?from={{.From}}">löschen</a>
							</li>
						{{else}}
							<li><i>keine Redirects</i></li>
						{{end}}
					</ul>
				</div>

				<div class="col-md-6">
					<h2>Neuer Kurzlink</h2>
					<form class="form-horizontal" action="/admin" method="post">
						<div class="form-group">
							<label for="from" class="col-sm-2 control-label">Kurzlink</label>
							<div class="col-sm-10">
								<input id="from" class="form-control" type="text" placeholder="z.B. lgt">
							</div>
						</div>

						<div class="form-group">
							<label for="to" class="col-sm-2 control-label">Zieladresse</label>
							<div class="col-sm-10">
								<input id="to" class="form-control" type="text" placeholder="z.B. http://www.yfu.de">
							</div>
						</div>

						<div class="form-group">
							<div class="col-sm-offset-2 col-sm-10">
								<button type="submit" class="btn btn-primary">Erstellen</button>
							</div>
						</div>
					</form>
				</div>
			</div>
		</div>
	</body>
</html>`
