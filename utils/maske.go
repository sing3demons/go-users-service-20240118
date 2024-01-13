package utils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

func MaskSensitiveData(data any) any {
	sensitiveFields := []string{"Password", "Email", "password", "email"}
	val := reflect.ValueOf(data)

	// Handle structs
	if val.Kind() == reflect.Struct {
		result := reflect.New(val.Type()).Elem()
		for i := 0; i < val.NumField(); i++ {
			fieldName := val.Type().Field(i).Name
			fieldValue := val.Field(i).Interface()
			// Check if the field is sensitive
			if contains(sensitiveFields, fieldName) {
				// Mask the sensitive value
				fieldValue = maskValue(fieldValue)
			}
			// Set the field in the result struct
			result.Field(i).Set(reflect.ValueOf(fieldValue))
		}

		return result.Interface()
	}

	// Handle maps
	if val.Kind() == reflect.Map {
		result := reflect.MakeMap(val.Type())

		// Iterate through map keys and values
		for _, key := range val.MapKeys() {
			fieldName := key.Interface().(string)
			fieldValue := val.MapIndex(key).Interface()
			// Check if the field is sensitive
			if contains(sensitiveFields, fieldName) {
				// Mask the sensitive value
				fieldValue = maskValue(fieldValue)
			}

			// Set the key-value pair in the result map
			result.SetMapIndex(key, reflect.ValueOf(fieldValue))
		}

		return result.Interface()
	}

	// Handle ptrs
	if val.Kind() == reflect.Ptr {
		// Recursively call MaskSensitiveData on the dereferenced value
		return MaskSensitiveData(val.Elem().Interface())
	}

	// Unsupported type, return as-is
	return data
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func maskValue(value interface{}) interface{} {
	// You can customize the masking logic based on the type of the value
	switch v := value.(type) {
	case string:
		return maskString(v)
	default:
		return value
	}
}

func maskString(input string) string {
	if isValidEmail(input) {
		return maskEmail(input)
	}

	return "********"
}

func ValidateBirthday(birthday string) (int64, error) {
	// birthday := "1997-05-26"
	parsedBirthday, err := time.Parse("2006-01-02", birthday)
	if err != nil {
		return 0, fmt.Errorf("error parsing birthday: %w", err)
	}

	timestamp := parsedBirthday.Unix()
	return timestamp, nil
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(emailRegex)
	return regex.MatchString(email)
}

func maskEmail(email string) string {
	parts := strings.Split(email, "@")
	localPart, domain := parts[0], parts[1]
	maskedLocalPart := string(localPart[0]) + strings.Repeat("*", len(localPart)-1)
	maskedEmail := fmt.Sprintf("%s@%s", maskedLocalPart, domain)

	return maskedEmail
}
