package media

import "strings"

// ParseCliOutput takes CLI output and returns a slice of strings.  Each item
// in the slice represents a single line of output
func ParseCliOutput(output string) []string {
	lines := make([]string, 0)
	retval := make([]string, 0)
	if len(strings.TrimSpace(output)) > 0 {
		lines = strings.Split(output, "\n")
	}

	//	Iterate through each entry.  If there is a blank line, remove it
	for _, item := range lines {
		if len(strings.TrimSpace(item)) > 0 {
			retval = append(retval, item)
		}
	}

	return retval
}
