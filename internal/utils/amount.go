package utils

import (
	"fmt"
	"strings"
)

// ParseAmountFlag parses a combined "number,currency" string (e.g. "60000,USD").
func ParseAmountFlag(s string) (number, currency string, err error) {
	parts := strings.SplitN(s, ",", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid amount %q: expected number,currency (e.g. 60000,USD)", s)
	}
	return parts[0], strings.ToUpper(strings.TrimSpace(parts[1])), nil
}
