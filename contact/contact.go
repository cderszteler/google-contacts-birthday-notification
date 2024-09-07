package contact

import (
	"context"
	"fmt"
	"google-contacts-birthday-notification/config"
	"google.golang.org/api/googleapi"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

type ServiceInterface interface {
	PersonFields(personFields string) ServiceInterface
	PageToken(pageToken string) ServiceInterface
	Do(opts ...googleapi.CallOption) (*people.ListConnectionsResponse, error)
}

type Service struct {
	ListConnections func(resourceName string) ServiceInterface
}

type Contact struct {
	Name     string
	Birthday people.Date
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

func NewContactService(config *config.Config) *Service {
	ctx := context.Background()
	credentials := parseCredentialsConfig(config)
	client := createClient(credentials, config)

	srv, err := people.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create people client %v", err)
	}
	return &Service{
		ListConnections: func(resourceName string) ServiceInterface {
			return &PeopleApiWrapper{call: srv.People.Connections.List(resourceName)}
		},
	}
}

func createClient(credentials *oauth2.Config, config *config.Config) *http.Client {
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

func parseCredentialsConfig(config *config.Config) *oauth2.Config {
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
