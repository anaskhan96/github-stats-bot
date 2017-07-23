package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"./analyze"
	"github.com/joho/godotenv"
	"github.com/jzelinskie/geddit"
)

// default interval for hitting the api in seconds
const interval = 2

func main() {
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	session, err := geddit.NewOAuthSession(env["clientID"], env["clientSecret"], env["userAgent"], env["redirectURL"])
	if err != nil {
		log.Fatal(err)
	}
	err = session.LoginAuth(env["username"], env["password"])
	if err != nil {
		log.Fatal(err)
	}

	go startPolling(session)

	// firing up a server to keep the process running
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Go program running %s!", r.URL.Path[1:])
}

func startPolling(s *geddit.OAuthSession) {
	tokenTime := time.Now()
	for {
		if time.Since(tokenTime).Hours() == 1 {
			s = refreshLogin()
			tokenTime = time.Now()
		}
		time.Sleep(interval * time.Second)
		go fetchComments(s)
	}
}

func fetchComments(s *geddit.OAuthSession) {
	comments, err := s.SubredditComments("all")
	if err != nil {
		log.Println(err)
		return
	}
	err = analyze.AnalyzeComments(s, comments)
	if err != nil {
		log.Println(err)
		return
	}
}

func refreshLogin() *geddit.OAuthSession {
	time.Sleep(interval * time.Second)
	log.Println("Authenticating a new login")
	var env map[string]string
	env, err := godotenv.Read()
	if err != nil {
		log.Fatal(err)
	}

	session, err := geddit.NewOAuthSession(env["clientID"], env["clientSecret"], env["userAgent"], env["redirectURL"])
	if err != nil {
		log.Fatal(err)
	}
	err = session.LoginAuth(env["username"], env["password"])
	if err != nil {
		log.Fatal(err)
	}
	return session
}
