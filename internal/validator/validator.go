package validator

import (
	"cmp"
	"fmt"
	"strings"
)

func IllegalArgument(field, desc string) error {
	return fmt.Errorf("invalid config %q %v", field, desc)
}

func Required(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return IllegalArgument(field, "is required")
	}
	return nil
}

func Coalesce[T cmp.Ordered](values ...T) T {
	var zero T
	for _, value := range values {
		// Test string type
		var item any = value
		if str, ok := item.(string); ok {
			if strings.TrimSpace(str) != "" {
				return value
			}
			continue
		}

		// Test other type
		if value > zero {
			return value
		}
	}
	return zero
}
