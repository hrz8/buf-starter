package main

import (
	"net/url"
)

// maskPassword masks the password in database URL for logging
func maskPassword(dbURL string) string {
	parsed, err := url.Parse(dbURL)
	if err != nil {
		return dbURL
	}

	if parsed.User != nil {
		if _, hasPassword := parsed.User.Password(); hasPassword {
			userInfo := parsed.User.Username() + ":***"
			maskedURL := parsed.Scheme + "://" + userInfo + "@" + parsed.Host + parsed.Path
			if parsed.RawQuery != "" {
				maskedURL += "?" + parsed.RawQuery
			}
			if parsed.Fragment != "" {
				maskedURL += "#" + parsed.Fragment
			}
			return maskedURL
		}
	}

	return parsed.String()
}
