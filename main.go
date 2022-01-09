package main

import (
	"context"
	"flag"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	listOpt := flag.Bool("l", false, "List aliases")
	aliasOpt := flag.String("a", "", "Add alias")
	flag.Parse()
	if flag.NArg() <= 0 {
		log.Fatalln("User key is required")
	}
	if !*listOpt && len(*aliasOpt) <= 0 {
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
	service, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalln(err)
	}

	if *listOpt {
		response, err := listAliases(service, key)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println("Current aliases.")
		for _, a := range response.Aliases {
			alias := a.(admin.Alias)
			fmt.Printf("%s - %s\n", alias.Alias, alias.PrimaryEmail)
		}
		return
	}

	response, err := addAlias(service, key, *aliasOpt)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Added new alias.")
	fmt.Println(response.Alias)
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

func listAliases(service *admin.Service, key string) (*admin.Aliases, error) {
	call := service.Users.Aliases.List(key)
	return call.Do()
}

func addAlias(service *admin.Service, key string, alias string) (*admin.Alias, error) {
	a := admin.Alias{
		Alias: alias,
	}
	call := service.Users.Aliases.Insert(key, &a)
	return call.Do()
}
