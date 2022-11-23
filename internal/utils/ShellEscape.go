package utils

import (
	"regexp"
	"strings"
)

var shellNoEscapeNeededRegexp = regexp.MustCompile("^[a-zA-Z0-9/_@=%:.,+-]*$")

// TODO: move elsewhere so it can be used for printing hints
func ShellEscape(str string) string {

	if shellNoEscapeNeededRegexp.MatchString(str) {
		return str
	}

	var out strings.Builder

	out.WriteRune('\'')

	for _, r := range str {

		switch r {
		case '\'':
			out.WriteString("'\\''")
		case '\n':
			out.WriteString("'\\n'")
		case '\r':
			out.WriteString("'\\r'")
		case '\t':
			out.WriteString("'\\t'")
		default:
			out.WriteRune(r)
		}
	}

	out.WriteRune('\'')

	return out.String()
}
