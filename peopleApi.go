package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var service *people.Service

func CreateService() {
	ctx := context.Background()
	credentials := createCredentialsConfig()
	client := createClient(credentials)

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people Client %v", err)
	}
	service = srv
}

func createClient(credentials *oauth2.Config) *http.Client {
	expiry, err := time.Parse(time.RFC3339, config.PeopleApi.Expiry)
	if err != nil {
		log.Fatalf("Unable to parse expiry %v", err)
	}
	return credentials.Client(context.Background(), &oauth2.Token{
		AccessToken:  config.PeopleApi.AccessToken,
		RefreshToken: config.PeopleApi.RefreshToken,
		TokenType:    config.PeopleApi.TokenType,
		Expiry:       expiry,
	})
}

func createCredentialsConfig() *oauth2.Config {
	credentials, err := os.ReadFile("credentials.local.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(credentials, people.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	return config
}
