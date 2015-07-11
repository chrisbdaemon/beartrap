/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package trap

import (
	"fmt"
	"strconv"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
	"github.com/chrisbdaemon/beartrap/trap/tcptrap"
)

// Interface defines the interface all traps adhere to
type Interface interface {
	Validate() []error
	Start() error
}

// Dispatcher provides interface for sending alerts to handlers
type Dispatcher interface {
	BroadcastAlert(alert.Alert)
}

// BaseTrap holds data common to all trap types
type BaseTrap struct {
	Severity int
	Dispatch Dispatcher
	params   config.Params
}

// New takes in a params object and returns a trap
func New(params config.Params, d Dispatcher) (Interface, error) {
	baseTrap := new(BaseTrap)
	var trap Interface

	baseTrap.Dispatch = d
	baseTrap.params = params

	switch params["type"] {
	case "tcp":
		trap = tcptrap.New(params, baseTrap)
	default:
		return nil, fmt.Errorf("Unknown trap type")
	}

	// will validate later *crosses fingers*
	baseTrap.Severity, _ = strconv.Atoi(params["severity"])

	return trap, nil
}

// TriggerAlert creates an alert and sends it to the broadcaster
func (trap *BaseTrap) TriggerAlert(s string) {
	a := alert.Alert{Message: s}
	trap.Dispatch.BroadcastAlert(a)
}

// Validate performs validation on the parameters of the trap
func (trap *BaseTrap) Validate() []error {
	errors := []error{}

	switch err := validate.Int(trap.params["severity"]); {
	case err != nil:
		errors = append(errors, fmt.Errorf("Invalid severity: %s", err))
	case trap.Severity < 0:
		errors = append(errors, fmt.Errorf("Severity cannot be negative"))
	}

	return errors
}
