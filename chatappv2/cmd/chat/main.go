package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"

	_ "github.com/joho/godotenv/autoload"
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
		return template.Must(template.ParseFiles(filepath.Join("cmd", "chat", "templates", t.filename)))
	}

	if os.Getenv("DEBUG") == "true" {
		t.template = loadCompileTemplate()
	} else {
		t.once.Do(func() {
			t.template = loadCompileTemplate()
		})
	}

	t.template.Execute(w, r)
}

func main() {
	// Parse addr flag
	addr := flag.String("addr", ":8080", "The server addr.")
	flag.Parse()

	// Setup OAuth2
	gomniauth.SetSecurityKey(os.Getenv("OAUTH_SECURITY_KEY"))
	gomniauth.WithProviders(
		google.New(os.Getenv("GOOGLE_OAUTH_CLIENT_ID"), os.Getenv("GOOGLE_OAUTH_SECRET"),
			fmt.Sprintf("%s://%s:%s/auth/callback/google", os.Getenv("PROTOCOL"), os.Getenv("HOST"), os.Getenv("PORT")),
		),
	)

	// Create and run room for chat
	r := newRoom()
	go r.run()

	// Assign handlers to paths
	http.Handle("/", MustAuth(&templateHandler{filename: "home.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Start listening
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
