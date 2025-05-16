package run

import (
	"fmt"
	"regexp"
	"strings"
)

var rfc1035Regex = regexp.MustCompile(`^[a-z]([-a-z0-9]{1,28}[a-z0-9])?$`)

func (c *Command) validate() []error {
	var errors []error

	if c.name == "" || !rfc1035Regex.MatchString(c.name) {
		errors = append(errors, fmt.Errorf("name must be provided and follow RFC1035 naming rules (3-30 characters, lowercase letters, digits, and hyphens, must start and end with a letter or digit)"))
	}

	if c.detatch && c.remove {
		errors = append(errors, fmt.Errorf("only one of detatch and remove can be provided at once"))
	}

	if c.detatch && c.interactive {
		errors = append(errors, fmt.Errorf("only one of detatch and interactive can be provided at once"))
	}

	if c.detatch && c.tty {
		errors = append(errors, fmt.Errorf("only one of detatch and tty can be provided at once"))
	}

	if strings.TrimSpace(c.image) == "" {
		errors = append(errors, fmt.Errorf("image must be provided"))
	}

	return errors
}
