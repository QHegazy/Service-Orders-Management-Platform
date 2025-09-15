package utils

import (
	"log"
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
		log.Printf("Email domain %s is a free email domain.\n", domain)
		return false
	}
	// mxRecords, err := net.LookupMX(domain)
	// if err != nil || len(mxRecords) == 0 {
	// 	return false
	// }

	return true
}
