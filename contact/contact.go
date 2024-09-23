package contact

import (
	"context"
	"fmt"
	"google-contacts-birthday-notification/config"
	"google.golang.org/api/googleapi"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type Service struct {
	ListConnections func(resourceName string) serviceInterface
}

type Contact struct {
	Name     string
	Birthday people.Date
}

type serviceInterface interface {
	PersonFields(personFields string) serviceInterface
	PageToken(pageToken string) serviceInterface
	Do(opts ...googleapi.CallOption) (*people.ListConnectionsResponse, error)
}

func (cs *Service) ListAllContacts() []Contact {
	var contacts = make([]Contact, 0)
	var nextPageToken string

	for {
		response, updated, err := cs.listContactsPaginated(nextPageToken, contacts)
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

func (cs *Service) listContactsPaginated(pageToken string, contacts []Contact) (*people.ListConnectionsResponse, []Contact, error) {
	response, err := cs.ListConnections("people/me").
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

func NewContactService(cfg config.Config) *Service {
	ctx := context.Background()
	credentials := parseCredentialsConfig(cfg)
	tokenSource := createRefreshedTokenSource(credentials, cfg)
	storeUpdatedToken(tokenSource, cfg)

	client := oauth2.NewClient(context.Background(), tokenSource)
	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people client %v", err)
	}
	return &Service{
		ListConnections: func(resourceName string) serviceInterface {
			return &peopleApiWrapper{call: srv.People.Connections.List(resourceName)}
		},
	}
}

func storeUpdatedToken(tokenSource oauth2.TokenSource, cfg config.Config) {
	token, err := tokenSource.Token()
	if err != nil {
		log.Fatalf("Unable to retrieve token: %v", err)
	}
	cfg.PeopleApi.Token = config.Token{
		AccessToken:  token.AccessToken,
		TokenType:    token.TokenType,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry.Format(time.RFC3339),
	}
	if err := config.WriteConfig(cfg); err != nil {
		log.Fatalf("Unable to write config file: %v", err)
	}
}

func createRefreshedTokenSource(credentials *oauth2.Config, config config.Config) oauth2.TokenSource {
	expiry, err := time.Parse(time.RFC3339, config.PeopleApi.Token.Expiry)
	if err != nil {
		log.Fatalf("Unable to parse expiry %v", err)
	}
	return credentials.TokenSource(context.Background(), &oauth2.Token{
		AccessToken:  config.PeopleApi.Token.AccessToken,
		RefreshToken: config.PeopleApi.Token.RefreshToken,
		TokenType:    config.PeopleApi.Token.TokenType,
		Expiry:       expiry,
	})
}

func parseCredentialsConfig(config config.Config) *oauth2.Config {
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
