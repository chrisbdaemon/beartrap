/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package handler

import (
	"fmt"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
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
	h        Interface
	receiver chan alert.Alert
	params   config.Params
}

// Start the underlying alert handler loop
func (handler *BaseHandler) Start() {
	for {
		a := <-handler.receiver
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

	return handler, nil
}

// Validate performs validation on the parameters of the handler
func (handler *BaseHandler) Validate() []error {
	errors := []error{}

	// Test any parameters here

	return errors
}
