package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/joho/godotenv"
	"github.com/jzelinskie/geddit"
)

// default interval for hitting the api in seconds
const interval = 3

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
	for {
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
	log.Println(comments[0].Body)
}
