package main

import (
	"fmt"
	"google-contacts-birthday-notification/config"
	"google-contacts-birthday-notification/contact"
	"google.golang.org/api/people/v1"
	"os"
	"strings"
	"time"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		processError(err)
	}

	mail := NewMailService(cfg)
	contacts := contact.NewContactService(cfg).ListAllContacts()
	tomorrow, soon := filterContacts(contacts)
	formatted := formatMail(tomorrow, soon)

	if err := mail.SendMail(formatted); err != nil {
		processError(err)
	}
}

func formatMail(tomorrow []contact.Contact, soon []contact.Contact) string {
	var builder strings.Builder
	builder.WriteString("Daily birthday notifications:")

	if len(tomorrow) > 0 {
		builder.WriteString("\n\nBirthday(s) tomorrow:")
		for _, contacts := range tomorrow {
			builder.WriteString(fmt.Sprintf("\n  - %s", contacts.Name))
		}
	}
	if len(soon) > 0 {
		builder.WriteString("\n\nBirthday(s) in two weeks:")
		for _, contacts := range soon {
			builder.WriteString(fmt.Sprintf("\n  - %s", contacts.Name))
		}
	}
	if len(tomorrow) == 0 && len(soon) == 0 {
		builder.WriteString("\n\nNo upcoming birthdays!")
	}

	return builder.String()
}

func filterContacts(contacts []contact.Contact) ([]contact.Contact, []contact.Contact) {
	tomorrowDate := time.Now().AddDate(0, 0, 1)
	soonDate := time.Now().AddDate(0, 0, 14)
	tomorrowFiltered := make([]contact.Contact, 0)
	soonFiltered := make([]contact.Contact, 0)

	for _, target := range contacts {
		if matchDateAndTime(target.Birthday, tomorrowDate) {
			tomorrowFiltered = append(tomorrowFiltered, target)
		} else if matchDateAndTime(target.Birthday, soonDate) {
			soonFiltered = append(soonFiltered, target)
		}
	}
	return tomorrowFiltered, soonFiltered
}

func matchDateAndTime(date people.Date, time time.Time) bool {
	return int(date.Day) == time.Day() &&
		int(date.Month) == int(time.Month())
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
