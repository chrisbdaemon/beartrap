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
	"strings"

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
	priority    syslog.Priority
	writer      *syslog.Writer
	host        string
	port        int
	proto       string
	tag         string
	params      config.Params
}

var syslogPriorityMap = map[string]syslog.Priority{
	"emergency": syslog.LOG_EMERG,
	"alert":     syslog.LOG_ALERT,
	"critical":  syslog.LOG_CRIT,
	"error":     syslog.LOG_ERR,
	"warning":   syslog.LOG_WARNING,
	"notice":    syslog.LOG_NOTICE,
	"info":      syslog.LOG_INFO,
	"debug":     syslog.LOG_DEBUG,
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

	sh.priority, _ = translateLogPriority(params["priority"])
	sh.params = params
	sh.port, _ = strconv.Atoi(params["port"])
	sh.host = params["host"]
	sh.proto = params["proto"]
	sh.tag = params["tag"]

	return sh
}

// Init sets up the logging functionality
func (sh *SyslogHandler) Init() error {
	var err error
	var writer *syslog.Writer

	if len(sh.host) > 0 {
		raddr := fmt.Sprintf("%s:%d", sh.host, sh.port)
		writer, err = syslog.Dial(sh.proto, raddr, sh.priority, sh.tag)
	} else {
		writer, err = syslog.New(sh.priority, sh.tag)
	}

	if err != nil {
		return fmt.Errorf("Unable to initialize syslog handler: %s", err)
	}

	sh.writer = writer
	return nil
}

// HandleAlert logs the alert text to the syslog daemon
func (sh *SyslogHandler) HandleAlert(a alert.Alert) {
	sh.writer.Write([]byte(a.Message))
}

func translateLogPriority(priorityParam string) (syslog.Priority, error) {
	priorities := make([]string, 0, 0)
	for k, v := range syslogPriorityMap {
		priorities = append(priorities, k)
		if k == priorityParam {
			return v, nil
		}
	}

	err := fmt.Errorf("Unknown log priority '%s', expecting: %s", priorityParam, strings.Join(priorities, ", "))
	return -1, err
}

// Validate performs validation on the trap and returns any errors
func (sh *SyslogHandler) Validate() []error {
	var errors []error

	errors = sh.baseHandler.Validate()

	// Priority is required and must be valid
	if priorityParam, ok := sh.params["priority"]; ok {
		_, err := translateLogPriority(priorityParam)
		if err != nil {
			errors = append(errors, err)
		}
	} else {
		errors = append(errors, fmt.Errorf("Missing syslog priority"))
	}

	// TODO: Ensure if one networking param is supplied, all others must be present
	if protoParam, ok := sh.params["proto"]; ok {
		if protoParam != "tcp" && protoParam != "udp" {
			errors = append(errors, fmt.Errorf("Invalid protocol '%s' expecting tcp or udp", protoParam))
		}
	}

	// Port is optional, if provided, make sure its valid
	if portParam, ok := sh.params["port"]; ok {
		err := validate.Port(portParam)
		if err != nil {
			errors = append(errors, fmt.Errorf("Invalid port: %s", err))
		}

		// If port is provided, make sure host is there as well
		if _, ok := sh.params["host"]; !ok {
			errors = append(errors, fmt.Errorf("syslog port provided but missing host"))
		}
	}

	// Host is optional, if provided, make sure its valid
	if hostParam, ok := sh.params["host"]; ok {
		err := validate.Host(hostParam)
		if err != nil {
			errors = append(errors, fmt.Errorf("Invalid host: %s", err))
		}
	}

	return errors
}
