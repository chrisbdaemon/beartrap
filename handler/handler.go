/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package handler

import (
	"fmt"
	"strconv"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
)

// Interface defines the interface all handlers adhere to
type Interface interface {
	Validate() []error
}

// BaseHandler holds data common to all handler types
type BaseHandler struct {
	Threshold int
	receiver  chan alert.Alert
	params    config.Params
}

// New takes in a params object and returns a handler
func New(params config.Params) (Interface, error) {
	baseHandler := new(BaseHandler)
	var handler Interface

	baseHandler.params = params

	switch params["type"] {
	case "syslog":
		// handler = syslog.New(params, baseHandler)
		handler = baseHandler
	default:
		return nil, fmt.Errorf("Unknown handler type")
	}

	// will validate later *crosses fingers*
	baseHandler.Threshold, _ = strconv.Atoi(params["threshold"])

	return handler, nil
}

// Validate performs validation on the parameters of the handler
func (handler *BaseHandler) Validate() []error {
	errors := []error{}

	switch err := validate.Int(handler.params["threshold"]); {
	case err != nil:
		errors = append(errors, fmt.Errorf("Invalid threshold: %s", err))
	case handler.Threshold < 0:
		errors = append(errors, fmt.Errorf("Threshold cannot be negative"))
	}

	return errors
}
