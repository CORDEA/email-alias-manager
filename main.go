package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	list := flag.Bool("l", false, "List aliases")
	alias := flag.String("a", "", "Add alias")
	flag.Parse()
	if flag.NArg() <= 0 {
		log.Fatalln("User key is required")
	}
	if !*list && len(*alias) <= 0 {
		log.Fatalln("Received illegal option")
	}
	key := flag.Arg(0)

	ctx := context.Background()
	path := os.Getenv("GOOGLE_CLIENT_SECRET")
	secret, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := authorize(ctx, secret)
	if err != nil {
		log.Fatalln(err)
	}
	_ = client
	_ = key
}

func authorize(ctx context.Context, secret []byte) (*http.Client, error) {
	config, err := google.ConfigFromJSON(secret, admin.AdminDirectoryUserAliasScope)
	if err != nil {
		return nil, err
	}
	state := generateState()
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	fmt.Println(url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return config.Client(ctx, token), nil
}

func generateState() string {
	state := ""
	for i := 0; i < 3; i++ {
		state += fmt.Sprintf("%c", rand.Intn(26)+97)
	}
	return state
}
