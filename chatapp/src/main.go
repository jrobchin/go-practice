package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

/*
	Create a template handler that compiles
	templates and puts it in the HTTP response.
*/
type templateHandler struct {
	once     sync.Once
	filename string
	template *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.template = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.template.Execute(w, nil)
}

func main() {
	r := newRoom()

	http.Handle("/", &templateHandler{filename: "home.html"})
	http.Handle("/room", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
