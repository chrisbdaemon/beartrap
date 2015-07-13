/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package sysloghandler

import (
	"fmt"
	"log"
	"log/syslog"
	"strconv"

	"github.com/chrisbdaemon/beartrap/alert"
	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/config/validate"
)

type handlerInterface interface {
	Validate() []error
	Start()
}

// SyslogHandler contains the details of a syslog handler
type SyslogHandler struct {
	baseHandler handlerInterface
	logger      *log.Logger
	host        string
	port        int
	params      config.Params
}

// Start starts up any handler-specific code then hands control over the the
// base handler's Start() method
func (sh *SyslogHandler) Start() {
	// Fire up any necessary connections.. hand reigns over to baseHandler
	log.Println("Firing base start()")
	sh.baseHandler.Start()
}

// New returns a new sysloghandler object waiting to handle alerts
func New(params config.Params, baseHandler handlerInterface) *SyslogHandler {
	sh := new(SyslogHandler)
	sh.baseHandler = baseHandler

	sh.params = params
	sh.port, _ = strconv.Atoi(params["port"])
	sh.host = params["host"]

	return sh
}

// Init sets up the logging functionality
// TODO: Allow priority and type to by set in config
func (sh *SyslogHandler) Init() error {
	logger, err := syslog.NewLogger(syslog.LOG_ALERT|syslog.LOG_AUTH, 12)
	if err != nil {
		return fmt.Errorf("Unable to initialize syslog handler: %s", err)
	}

	logger.SetFlags(0)
	sh.logger = logger
	return nil
}

// HandleAlert logs the alert text to the syslog daemon
func (sh *SyslogHandler) HandleAlert(a alert.Alert) {
	sh.logger.Println(a.Message)
}

// Validate performs validation on the trap and returns any errors
func (sh *SyslogHandler) Validate() []error {
	var errors []error

	errors = sh.baseHandler.Validate()

	if portParam, ok := sh.params["port"]; ok {
		err := validate.Port(portParam)
		if err != nil {
			errors = append(errors, fmt.Errorf("Invalid port: %s", err))
		}
	}

	if hostParam, ok := sh.params["host"]; ok {
		err := validate.Host(hostParam)
		if err != nil {
			errors = append(errors, fmt.Errorf("Invalid host: %s", err))
		}
	}

	return errors
}
