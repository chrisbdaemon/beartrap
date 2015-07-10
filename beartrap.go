/*
 * Copyright (c) 2015, Chris Benedict <chrisbdaemon@gmail.com>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * Redistributions of source code must retain the above copyright notice, this
 * list of conditions and the following disclaimer.
 *
 * Redistributions in binary form must reproduce the above copyright notice,
 * this list of conditions and the following disclaimer in the documentation
 * and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chrisbdaemon/beartrap/config"
	"github.com/chrisbdaemon/beartrap/trap"
	getopt "github.com/kesselborn/go-getopt"
)

func main() {
	options := getOptions()

	cfg, err := config.New(options["config"].String)
	if err != nil {
		log.Fatal(err)
	}

	trapParams, err := cfg.TrapParams()
	if err != nil {
		log.Fatalf("Error reading traps: %s", err)
	}

	traps, err := initTraps(trapParams)
	if err != nil {
		log.Fatalf("Error initializing traps: %s", err)
	}
	errors := validateTraps(traps)

	if len(errors) > 0 {
		displayErrors(errors)
		os.Exit(-1)
	} else {
		startTraps(traps)

		// Hack to let traps run till I build out trap ctrl more
		for {
			time.Sleep(500 * time.Second)
		}
	}
}

func validateTraps(traps []trap.Interface) []error {
	var errors []error
	for i := range traps {
		errors = append(errors, traps[i].Validate()...)
	}
	return errors
}

// displayErrors takes a slice of errors and prints them to the screen
func displayErrors(errors []error) {
	for i := range errors {
		log.Println(errors[i])
	}
}

// startTraps takes a slice of traps and starts them in a goroutine
// TODO: Allow them to be stopped
func startTraps(traps []trap.Interface) {
	for i := range traps {
		go traps[i].Start()
	}
}

// initTraps take in a list of trap parameters, creates trap objects
// that are returned along with any errors generated from validation
func initTraps(trapParams []config.Params) ([]trap.Interface, error) {
	traps := []trap.Interface{}

	for i := range trapParams {
		trap, err := trap.New(trapParams[i])
		if err != nil {
			return nil, err
		}

		traps = append(traps, trap)
	}

	return traps, nil
}

// Parse commandline arguments into getopt object
func getOptions() map[string]getopt.OptionValue {
	optionDefinition := getopt.Options{
		Description: "Beartrap v0.3 by Chris Benedict <chrisbdaemon@gmail.com>",
		Definitions: getopt.Definitions{
			{"config|c|BEARTRAP_CONFIG", "configuration file", getopt.Required, ""},
		},
	}

	options, _, _, err := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]

	if err != nil || wantsHelp {
		exitCode := 0

		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", err.Error(), "\n", optionDefinition.Help())
			exitCode = err.ErrorCode
		}
		os.Exit(exitCode)
	}

	return options
}
