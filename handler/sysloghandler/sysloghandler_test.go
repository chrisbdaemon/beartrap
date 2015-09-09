/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Licensing terms are located in LICENSE file.
 */

package sysloghandler

import (
	"fmt"
	"log/syslog"
	"math/rand"
	"net"
	"testing"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/stretchr/testify/assert"
)

type stubBaseHandler struct{}

func (h stubBaseHandler) Start() {}
func (h stubBaseHandler) Validate() []error {
	var errors []error
	return errors
}

func TestNew(t *testing.T) {
	var baseHandler stubBaseHandler
	params := config.Params{}

	handler := New(params, baseHandler)
	assert.NotNil(t, handler)
}

var baseParams = config.Params{
	"proto":    "tcp",
	"priority": "alert",
	"host":     "8.8.8.8",
	"port":     "12345",
}

func TestValidate(t *testing.T) {
	bh := &stubBaseHandler{}
	params := baseParams

	sh := New(params, bh)
	errors := sh.Validate()
	assert.Len(t, errors, 0)

	params["proto"] = "foobar"
	params["host"] = "foobar"
	params["port"] = "foobar"
	params["priority"] = "panic"
	sh = New(params, bh)
	errors = sh.Validate()
	assert.Len(t, errors, 4)
}

func TestInit(t *testing.T) {
	sh := &SyslogHandler{}
	sh.proto = "tcp"
	sh.priority = syslog.LOG_ALERT
	sh.host = "127.0.0.1"
	sh.port = 1024 + (rand.Int() % 64511)

	fakeServer("tcp", sh.port)

	err := sh.Init()
	assert.NoError(t, err)
}

func fakeServer(proto string, port int) {
	raddr := fmt.Sprintf(":%d", port)
	ln, _ := net.Listen(proto, raddr)
	go func(l net.Listener) {
		conn, _ := ln.Accept()
		conn.Close()
	}(ln)
}
