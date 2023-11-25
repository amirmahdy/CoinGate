package utils

import (
	"fmt"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\\s]+$`).MatchString
)

func ValidateString(str string, min int, max int) error {
	if len(str) < min || len(str) > max {
		return fmt.Errorf("string length must be between %d and %d", min, max)
	}
	return nil
}

func ValidateUsername(name string) error {
	err := ValidateString(name, 3, 50)
	if err != nil {
		return err
	}
	if !isValidUsername(name) {
		return fmt.Errorf("username must contain only letters, numbers and underscores")
	}
	return nil
}

func ValidateFullname(name string) error {
	if err := ValidateString(name, 3, 50); err != nil {
		return err
	}
	if !isValidFullname(name) {
		return fmt.Errorf("fullname must contain only letters and spaces")
	}
	return nil
}
