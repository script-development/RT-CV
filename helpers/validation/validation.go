package validation

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidDomain is the error message send if the domain name is invalid
var ErrInvalidDomain = errors.New("invalid domain name")

// ValidDomain returns an error if the domain is not valid.
func ValidDomain(domain string, allowWildcard bool) error {
	parts := strings.Split(domain, ".")
	for _, part := range parts {
		if part == "" {
			return ErrInvalidDomain
		}
		if part == "*" {
			if allowWildcard {
				continue
			}
			return ErrInvalidDomain
		}
		if len(parts) == 1 {
			// Only wildcard match all domain names are allowed
			return ErrInvalidDomain
		}
		for idx, letter := range part {
			if (letter >= 'a' && letter <= 'z') || (letter >= 'A' && letter <= 'Z') || (letter >= '0' && letter <= '9') {
				continue
			}
			if letter == '-' && idx != 0 && idx != len(part)-1 {
				// "-" is not allowed as first and last character of a domain name
				continue
			}
			return ErrInvalidDomain
		}
	}
	return nil
}

// ValidDomainListAndFormat formats the domains list and check if there are invalid domain names
func ValidDomainListAndFormat(domains *[]string, allowWildcard bool) error {
	for idx, domain := range *domains {
		domain = strings.TrimSpace(strings.ToLower(domain))
		if ValidDomain(domain, allowWildcard) != nil {
			return fmt.Errorf("domain %d is invalid", idx)
		}
		(*domains)[idx] = domain
	}
	return nil
}
