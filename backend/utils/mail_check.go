package utils

import (
	"fmt"
	"regexp"
	"strings"
)

func IsBusinessEmail(email string) bool {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]

	freeEmailDomainsRegex := regexp.MustCompile(
		`^(?i)(gmail|yahoo|hotmail|outlook|aol|icloud|protonmail|mail|zoho|yandex|gmx|live|msn)\.(com|net|org|co\.uk|fr|de|ru)$`,
	)

	if freeEmailDomainsRegex.MatchString(domain) {
		return false
	}

	// mxRecords, err := net.LookupMX(domain)
	// if err != nil || len(mxRecords) == 0 {
	// 	return false
	// }

	return true
}

func MockSendMail(to string, subject string, body string) error {
	fmt.Printf("Mocking email send to %s with subject %s and body %s\n", to, subject, body)
	return nil
}
