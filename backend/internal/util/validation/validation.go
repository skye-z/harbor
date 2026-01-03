package validation

import (
	"regexp"
	"strings"
)

func ValidateRequired(value, fieldName string) bool {
	return strings.TrimSpace(value) != ""
}

func ValidateID(id string) bool {
	return len(id) > 0
}

func ValidatePath(path string) bool {
	if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
		return false
	}
	return true
}

func ValidateTag(tag string) bool {
	return strings.TrimSpace(tag) != ""
}

func ValidateName(name string) bool {
	if len(name) == 0 || len(name) > 255 {
		return false
	}
	return true
}

func ValidateDockerTag(tag string) bool {
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9_.-]*(:[\w][\w.-]*)?$`
	matched, _ := regexp.MatchString(pattern, tag)
	return matched
}

func ValidateUsername(username string) bool {
	if len(username) < 3 || len(username) > 32 {
		return false
	}
	pattern := `^[a-zA-Z0-9_]+$`
	matched, _ := regexp.MatchString(pattern, username)
	return matched
}

func ValidatePassword(password string) bool {
	return len(password) >= 8
}
