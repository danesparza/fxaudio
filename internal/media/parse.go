package media

import "strings"

// ParseCliOutput takes CLI output and returns a slice of strings.  Each item
// in the slice represents a single line of output
func ParseCliOutput(output string) []string {
	retval := make([]string, 0)

	if len(strings.TrimSpace(output)) > 0 {
		retval = strings.Split(output, "\n")
	}

	return retval
}
