package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/sirupsen/logrus"
)

type provider struct {
	name     string
	id       string
	secret   string
	callback string
	class2   string
	whatever string
}

func main() {
	logrus.Printf("Hello from the Log")
	err := godotenv.Load()
	if err != nil {
		logrus.Printf("No .env file, resorting to environment")
	}
	key := os.Getenv("SESSION_KEY")
	if key == "" {
		logrus.Fatal("Environment variables missing, can't continue.")
	}

	clientID := os.Getenv("GOOGLE_OAUTH2_CLIENT_ID")
	if clientID == "" {
		logrus.Fatalf("Google client id missing")
	}
	clientSecret := os.Getenv("GOOGLE_OAUTH2_CLIENT_SECRET")
	if clientSecret == "" {
		logrus.Fatalf("Google client secret missing")
	}

	maxAge := 86400 * 30 // Number of days auth valid for
	isProd := false

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(clientID,
			clientSecret,
			"http://lvh.me:3000/auth/google/callback",
			"email",
			"profile"),
	)

	p := pat.New()
	p.Get("/auth/{provider}/callback", func(res http.ResponseWriter, req *http.Request) {
		user, err := gothic.CompleteUserAuth(res, req)
		if err != nil {
			fmt.Fprintln(res, err)
			return
		}
		t, _ := template.ParseFiles("templates/success.html")
		t.Execute(res, user)
	})

	p.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.BeginAuthHandler(res, req)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")
		t.Execute(res, false)
	})
	logrus.Println("listening on localhost:3000")
	logrus.Fatal(http.ListenAndServe(":3000", p))
}
