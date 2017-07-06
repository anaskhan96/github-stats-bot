package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/jzelinskie/geddit"
)

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
	log.Println(session.Me())
}
