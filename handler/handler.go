/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
	"github.com/chrisbdaemon/beartrap/handler/sysloghandler"
)

// Interface defines the interface all handlers adhere to
type Interface interface {
	Validate() []error
	Start()
	HandleAlert(alert.Alert)
	Init() error
}

// BaseHandler holds data common to all handler types
type BaseHandler struct {
	Threshold int
	h         Interface
	receiver  chan alert.Alert
	params    config.Params
}

// Start the underlying alert handler loop
func (handler *BaseHandler) Start() {
	for {
		log.Println("Checking for alerts")
		a := <-handler.receiver
		log.Println("Got one")
		handler.h.HandleAlert(a)
	}
}

// New takes in a params object and returns a handler
func New(params config.Params, alertChan chan alert.Alert) (Interface, error) {
	baseHandler := new(BaseHandler)
	var handler Interface

	baseHandler.params = params
	baseHandler.receiver = alertChan

	switch params["type"] {
	case "syslog":
		handler = sysloghandler.New(params, baseHandler)
	default:
		return nil, fmt.Errorf("Unknown handler type")
	}

	baseHandler.h = handler

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
