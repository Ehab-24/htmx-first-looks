package lib

import (
	"strings"
)

func KebabCase(str string) string {
	return strings.ReplaceAll(strings.ToLower(str), " ", "-")
}
