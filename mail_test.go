package main

import (
	"errors"
	"testing"

	"github.com/wneessen/go-mail"
)

var mockConfig = Config{
	Mail: Mail{
		Host:     "smtp.example.com",
		Port:     587,
		User:     "testuser",
		Password: "testpassword",
		Tls:      true,
		Secure:   true,
		Sender:   "sender@example.com",
		Receiver: "receiver@example.com",
	},
}

// MockClient is a mock mail client for testing
type MockClient struct {
	DialAndSendCalled bool
	DialAndSendError  error
}

func (m *MockClient) DialAndSend(_ ...*mail.Msg) error {
	m.DialAndSendCalled = true
	return m.DialAndSendError
}

func TestCreateMailClient(t *testing.T) {
	origConfig := config
	origClient := client
	defer func() {
		config = origConfig
		client = origClient
	}()
	config = mockConfig

	// Test TLSMandatory
	config.Mail.Tls = true
	config.Mail.Secure = true
	testCreation(t)

	// Test TLSOpportunistic
	config.Mail.Tls = true
	config.Mail.Secure = false
	testCreation(t)

	// Test NoTLS
	config.Mail.Tls = false
	config.Mail.Secure = false
	testCreation(t)
}

func testCreation(t *testing.T) {
	createMailClient()
	if client == nil {
		t.Error("Expected client to be created, but it's nil")
	}
}

func TestSendMail(t *testing.T) {
	origConfig := config
	origClient := client
	defer func() {
		config = origConfig
		client = origClient
	}()
	config = mockConfig

	// Create a mock client
	mockClient := &MockClient{}
	client = mockClient

	// Test successful send
	err := SendMail("Test message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !mockClient.DialAndSendCalled {
		t.Error("Expected DialAndSend to be called, but it wasn't")
	}

	// Test error in DialAndSend
	expectedError := errors.New("failed to send email")
	mockClient.DialAndSendError = expectedError
	err = SendMail("Test message")
	if !errors.Is(expectedError, err) {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	// Test error in setting From address
	config.Mail.Sender = "invalid-email"
	err = SendMail("Test message")
	if err == nil {
		t.Error("Expected an error when setting an invalid From address, but got none")
	}

	// Test error in setting To address
	config.Mail.Sender = "sender@example.com"
	config.Mail.Receiver = "invalid-email"
	err = SendMail("Test message")
	if err == nil {
		t.Error("Expected an error when setting an invalid To address, but got none")
	}
}
