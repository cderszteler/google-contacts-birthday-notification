package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

var service *people.Service

func CreateService() {
	ctx := context.Background()
	credentials := parseCredentialsConfig()
	client := createClient(credentials)

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people client %v", err)
	}
	service = srv
}

func createClient(credentials *oauth2.Config) *http.Client {
	expiry, err := time.Parse(time.RFC3339, config.PeopleApi.Token.Expiry)
	if err != nil {
		log.Fatalf("Unable to parse expiry %v", err)
	}
	return credentials.Client(context.Background(), &oauth2.Token{
		AccessToken:  config.PeopleApi.Token.AccessToken,
		RefreshToken: config.PeopleApi.Token.RefreshToken,
		TokenType:    config.PeopleApi.Token.TokenType,
		Expiry:       expiry,
	})
}

func parseCredentialsConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.PeopleApi.Credentials.ClientId,
		ClientSecret: config.PeopleApi.Credentials.ClientSecret,
		RedirectURL:  config.PeopleApi.Credentials.RedirectUris[0],
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.PeopleApi.Credentials.AuthUri,
			TokenURL: config.PeopleApi.Credentials.TokenUri,
		},
	}
}
