package contact

import (
	"google-contacts-birthday-notification/config"
	"google.golang.org/api/googleapi"
	"testing"
	"time"

	"google.golang.org/api/people/v1"
)

var mockContactService = &Service{
	ListConnections: func(resourceName string) serviceInterface {
		return &mockPeopleConnectionsListCall{}
	},
}

type mockPeopleConnectionsListCall struct {
	nextPageToken string
}

func (m *mockPeopleConnectionsListCall) PageToken(pageToken string) serviceInterface {
	return m
}

func (m *mockPeopleConnectionsListCall) PersonFields(personFields string) serviceInterface {
	return m
}

func (m *mockPeopleConnectionsListCall) Do(opts ...googleapi.CallOption) (*people.ListConnectionsResponse, error) {
	return &people.ListConnectionsResponse{
		Connections: []*people.Person{
			{
				Names: []*people.Name{
					{DisplayName: "John Doe"},
				},
				Birthdays: []*people.Birthday{
					{Date: &people.Date{Year: 1990, Month: 1, Day: 1}},
				},
			},
			{
				Names: []*people.Name{
					{DisplayName: "Jane Doe"},
				},
				Birthdays: []*people.Birthday{
					{Date: &people.Date{Year: 1992, Month: 2, Day: 2}},
				},
			},
		},
	}, nil
}

func TestListAllContacts(t *testing.T) {
	contacts := mockContactService.ListAllContacts()

	if len(contacts) != 2 {
		t.Errorf("Expected 2 contacts, got %d", len(contacts))
	}

	expectedNames := []string{"John Doe", "Jane Doe"}
	for i, contact := range contacts {
		if contact.Name != expectedNames[i] {
			t.Errorf("Expected name %s, got %s", expectedNames[i], contact.Name)
		}
	}
}

func TestCreateService(t *testing.T) {
	mockConfig := config.Config{
		PeopleApi: config.PeopleApi{
			Token: config.Token{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				TokenType:    "Bearer",
				Expiry:       time.Now().Add(time.Hour).Format(time.RFC3339),
			},
			Credentials: config.Credentials{
				ClientId:     "client_id",
				ClientSecret: "client_secret",
				RedirectUris: []string{"http://localhost"},
				AuthUri:      "https://accounts.google.com/o/oauth2/auth",
				TokenUri:     "https://oauth2.googleapis.com/token",
			},
		},
	}
	service := NewContactService(mockConfig)

	if service == nil {
		t.Error("Expected service to be created, but it's nil")
	}
}

func TestParseCredentialsFromConfig(t *testing.T) {
	mockConfig := config.Config{
		PeopleApi: config.PeopleApi{
			Credentials: config.Credentials{
				ClientId:     "client_id",
				ClientSecret: "client_secret",
				RedirectUris: []string{"http://localhost"},
				AuthUri:      "https://accounts.google.com/o/oauth2/auth",
				TokenUri:     "https://oauth2.googleapis.com/token",
			},
		},
	}
	credentials := parseCredentialsConfig(mockConfig)

	if credentials.ClientID != "client_id" {
		t.Errorf("Expected ClientID 'client_id', got %s", credentials.ClientID)
	}
	if credentials.ClientSecret != "client_secret" {
		t.Errorf("Expected ClientSecret 'client_secret', got %s", credentials.ClientSecret)
	}
	if credentials.RedirectURL != "http://localhost" {
		t.Errorf("Expected RedirectURL 'http://localhost', got %s", credentials.RedirectURL)
	}
}
