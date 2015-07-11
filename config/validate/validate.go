/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package validate

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Int takes in a value, ensures it is not empty and contains a valid integer.
// if either of the checks fails it returns an error
func Int(value string) error {
	// only checking for error
	_, err := strconv.Atoi(value)

	switch {
	case len(strings.TrimSpace(value)) < 1:
		return fmt.Errorf("missing value")
	case err != nil:
		return fmt.Errorf("expecting integer, got '%s'", value)
	}

	return nil
}

// Port ensures first the value is a valid integer, if so it checks if it
// is within the valid port range (1-65535). If it is not valid, returns an error
func Port(value string) error {
	// make sure its a valid integer first
	if err := Int(value); err != nil {
		return err
	}

	// already know its good, ignore err
	i, _ := strconv.Atoi(value)

	if i < 1 || i > 65535 {
		return fmt.Errorf("%s not within range of valid ports (1-65535)", value)
	}

	return nil
}

// Host validates a host to ensure it is valid. If not, it returns an error
func Host(value string) error {
	// only care about err
	_, err := net.ResolveIPAddr("ip", value)

	if len(strings.TrimSpace(value)) < 1 {
		return fmt.Errorf("missing host")
	}

	if err != nil {
		return fmt.Errorf("invalid host, expecting IP address got '%s'", value)
	}

	return nil
}
