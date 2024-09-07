package main

import (
	"errors"
	"google-contacts-birthday-notification/config"
	"testing"

	"github.com/wneessen/go-mail"
)

var mockConfig = &config.Config{
	Mail: config.Mail{
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
type mockClient struct {
	DialAndSendCalled bool
	DialAndSendError  error
}

func (m *mockClient) DialAndSend(_ ...*mail.Msg) error {
	m.DialAndSendCalled = true
	return m.DialAndSendError
}

func TestCreateMailClient(t *testing.T) {
	origConfig := mockConfig
	defer func() {
		mockConfig = origConfig
	}()

	// Test TLSMandatory
	mockConfig.Mail.Tls = true
	mockConfig.Mail.Secure = true
	testCreation(t)

	// Test TLSOpportunistic
	mockConfig.Mail.Tls = true
	mockConfig.Mail.Secure = false
	testCreation(t)

	// Test NoTLS
	mockConfig.Mail.Tls = false
	mockConfig.Mail.Secure = false
	testCreation(t)
}

func testCreation(t *testing.T) {
	client := NewMailService(mockConfig)
	if client == nil {
		t.Error("Expected client to be created, but it's nil")
	}
}

func TestSendMail(t *testing.T) {
	client := &mockClient{}
	service := &MailService{
		config: mockConfig,
		client: client,
	}

	// Test successful send
	err := service.SendMail("Test message")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !client.DialAndSendCalled {
		t.Error("Expected DialAndSend to be called, but it wasn't")
	}

	// Test error in DialAndSend
	expectedError := errors.New("failed to send email")
	client.DialAndSendError = expectedError
	err = service.SendMail("Test message")
	if !errors.Is(expectedError, err) {
		t.Errorf("Expected %v, got %v", expectedError, err)
	}

	// Test error in setting From address
	mockConfig.Mail.Sender = "invalid-email"
	err = service.SendMail("Test message")
	if err == nil {
		t.Error("Expected an error when setting an invalid From address, but got none")
	}

	// Test error in setting To address
	mockConfig.Mail.Sender = "sender@example.com"
	mockConfig.Mail.Receiver = "invalid-email"
	err = service.SendMail("Test message")
	if err == nil {
		t.Error("Expected an error when setting an invalid To address, but got none")
	}
}
