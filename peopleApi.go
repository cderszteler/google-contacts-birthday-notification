package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type Contact struct {
	Name     string
	Birthday people.Date
}

var service *people.Service

func ListAllContacts() []Contact {
	var contacts = make([]Contact, 0)
	var nextPageToken string
	if service == nil {
		createService()
	}

	for {
		response, updated, err := listContactsPaginated(nextPageToken, contacts)
		if err != nil {
			log.Fatalf("Error listing contacts: %v", err)
		}

		contacts = updated
		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return contacts
}

func listContactsPaginated(pageToken string, contacts []Contact) (*people.ListConnectionsResponse, []Contact, error) {
	response, err := service.People.Connections.List("people/me").
		PageToken(pageToken).
		PersonFields("names,birthdays").
		Do()

	if err != nil {
		return nil, contacts, fmt.Errorf("unable to retrieve people: %v", err)
	}

	for _, connection := range response.Connections {
		if len(connection.Names) > 0 && len(connection.Birthdays) > 0 {
			name := connection.Names[0].DisplayName
			birthday := connection.Birthdays[0].Date
			contacts = append(contacts, Contact{Name: name, Birthday: *birthday})
		}
	}
	return response, contacts, nil
}

func createService() {
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
