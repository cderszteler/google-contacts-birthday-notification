package scripts

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/people/v1"
	"log"
	"os"
)

var config *oauth2.Config

// generateToken This function can be called manually to authenticate with Google via OAuth
// and print a token pair that is required for the configuration file.
//
//goland:noinspection GoUnusedFunction
func generateToken() {
	createConfig()
	code := generateAuthUrl()
	decodeCode(code)
}

func createConfig() {
	credentials, err := os.ReadFile("credentials.local.json")
	config, err = google.ConfigFromJSON(credentials, people.ContactsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
}

func generateAuthUrl() string {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code (query parameter): \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	return authCode
}

func decodeCode(code string) {
	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	fmt.Printf("%+v\n", tok)
}
