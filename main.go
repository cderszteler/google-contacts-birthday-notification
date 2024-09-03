package main

import (
	"fmt"
	"google.golang.org/api/people/v1"
	"os"
	"strings"
	"time"
)

var config Config

func main() {
	if err := ReadConfig(&config); err != nil {
		processError(err)
	}
	contacts := ListAllContacts()
	tomorrow, soon := filterContacts(contacts)
	formatted := formatMail(tomorrow, soon)

	if err := SendMail(formatted); err != nil {
		processError(err)
	}
}

func formatMail(tomorrow []Contact, soon []Contact) string {
	var builder strings.Builder
	builder.WriteString("Daily birthday notifications:")

	if len(tomorrow) > 0 {
		builder.WriteString("\n\nBirthday(s) tomorrow:")
		for _, contact := range tomorrow {
			builder.WriteString(fmt.Sprintf("\n  - %s", contact.Name))
		}
	}
	if len(soon) > 0 {
		builder.WriteString("\n\nBirthday(s) in two weeks:")
		for _, contact := range soon {
			builder.WriteString(fmt.Sprintf("\n  - %s", contact.Name))
		}
	}
	if len(tomorrow) == 0 && len(soon) == 0 {
		builder.WriteString("\n\nNo upcoming birthdays!")
	}

	return builder.String()
}

func filterContacts(contacts []Contact) ([]Contact, []Contact) {
	tomorrowDate := time.Now().AddDate(0, 0, 1)
	soonDate := time.Now().AddDate(0, 0, 14)
	tomorrowFiltered := make([]Contact, 0)
	soonFiltered := make([]Contact, 0)

	for _, contact := range contacts {
		if matchDateAndTime(contact.Birthday, tomorrowDate) {
			tomorrowFiltered = append(tomorrowFiltered, contact)
		} else if matchDateAndTime(contact.Birthday, soonDate) {
			soonFiltered = append(soonFiltered, contact)
		}
	}
	return tomorrowFiltered, soonFiltered
}

func matchDateAndTime(date people.Date, time time.Time) bool {
	return int(date.Day) == time.Day() &&
		int(date.Month) == int(time.Month()) &&
		(date.Year == 0 || int(date.Year) == time.Year())
}

func processError(err error) {
	fmt.Println(err)
	os.Exit(2)
}
