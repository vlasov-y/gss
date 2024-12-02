package config

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strings"
)

func DecodeACMEEmail(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(ACMEEmail("")) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("invalid email address: ACMEEmail expects a string, got %T", data)
	}
	value := data.(string)
	// Email validation logic
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, value)
	if !match {
		return nil, fmt.Errorf("invalid email address: %s", value)
	}
	return ACMEEmail(value), nil
}

func DecodeACMEURL(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(ACMEURL("")) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("invalid ACME URL: ACMEURL expects a string, got %T", data)
	}
	value := data.(string)
	// URL validation logic
	parsedURL, err := url.ParseRequestURI(value)
	if err != nil {
		return nil, fmt.Errorf("invalid ACME URL: %w", err)
	}
	return ACMEURL(parsedURL.String()), nil
}

func DecodeACMEDomains(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(ACMEDomains{}) {
		return data, nil
	}
	parseDomains := func(domains []string) (ACMEDomains, error) {
		domainRegex := `^(\*\.)?([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`
		var result ACMEDomains
		for _, domain := range domains {
			domain = strings.TrimSpace(domain)
			if !regexp.MustCompile(domainRegex).MatchString(domain) {
				return nil, fmt.Errorf("invalid domain name: %s", domain)
			}
			if slices.Contains(result, domain) {
				return nil, fmt.Errorf("duplicate domain name: %s", domain)
			}
			result = append(result, domain)
		}
		return result, nil
	}
	// If the input is a string, split it into a slice of domains
	if f.Kind() == reflect.String {
		domainsStr := data.(string)
		// Handle comma-separated domains
		domains := strings.Split(domainsStr, ",")
		return parseDomains(domains)
	}
	// If the input is already a []string, parse it directly
	if f.Kind() == reflect.Slice {
		domains := []string{}
		slice, ok := data.([]interface{})
		if !ok {
			return nil, errors.New("unsupported type for ACME domains: cannot cast slice to []interface{}")
		}
		for i, object := range slice {
			s, ok := object.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported type for ACME domains at index %d, expected string, got %s", i, reflect.TypeOf(object))
			}
			domains = append(domains, s)
		}
		return parseDomains(domains)
	}
	// If not a string or []string, return an error
	return nil, fmt.Errorf("unsupported type for ACME domains, expected string or []string, got %s", f.Kind())
}

func DecodeACMEChallengePath(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t != reflect.TypeOf(ACMEChallengePath("")) {
		return data, nil
	}
	if f != reflect.TypeFor[string]() {
		return nil, fmt.Errorf("invalid ACME challenge path: ACMEChallengePath expects a string, got %T", data)
	}
	value := data.(string)
	// Path validation logic
	if !strings.HasPrefix(value, "/") {
		return nil, fmt.Errorf("invalid ACME challenge path: ACME challenge path must start with '/': %s", value)
	}
	parsedURL, err := url.Parse(value)
	if err != nil || parsedURL.Path != value {
		return nil, fmt.Errorf("invalid ACME challenge path: %s", value)
	}
	return ACMEChallengePath(value), nil
}
