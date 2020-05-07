package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/mgo.v2"
)

type contextKey struct {
	name string
}

var contextKeyAPIKey = &contextKey{"api-key"}

// APIKey extracts the API key from context as well as an
// ok bool indicating whether or not the key was successful extracted
func APIKey(ctx context.Context) (string, bool) {
	key, ok := ctx.Value(contextKeyAPIKey).(string)
	return key, ok
}

func isValidAPIKey(key string) bool {
	return key == "wowsosecret"
}

func withAPIKey(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if !isValidAPIKey(key) {
			respondErr(w, r, http.StatusUnauthorized, "invalid API key")
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyAPIKey, key)
		fn(w, r.WithContext(ctx))
	}
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}

// Server is the API server
type Server struct {
	db *mgo.Session
}

func main() {
	err := godotenv.Load(os.Getenv("DOTENV"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		addr  = flag.String("addr", ":8080", "endpoint address")
		mongo = flag.String("mongo", os.Getenv("MONGO_URI"), "mongodb address")
	)
	flag.Parse()

	log.Println("Dialing mongo", *mongo)
	db, err := mgo.Dial(*mongo)
	if err != nil {
		log.Fatalln("failed to connect to mongo:", err)
	}
	defer db.Close()

	s := &Server{
		db: db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/polls/", withCORS(withAPIKey(s.handlePolls)))

	log.Println("Starting web server on", *addr)
	http.ListenAndServe(*addr, mux)
	log.Println("Stopping...")
}
