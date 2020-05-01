package main

import (
	"log"
	"net/http"
	"os"
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
	loadCompileTemplate := func() *template.Template {
		return template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	}

	if os.Getenv("DEBUG") == "true" {
		t.template = loadCompileTemplate()
	} else {
		t.once.Do(func() {
			t.template = loadCompileTemplate()
		})
	}

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
