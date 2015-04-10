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

	"bitbucket.com/chrisbdaemon/beartrap/config"
	"github.com/kesselborn/go-getopt"
)

func main() {
	options := getOptions()
	cfg, err := config.New(options["config"].String)

	if err != nil {
		log.Fatal(err)
	}

	cfg.TrapParams()
}

func getOptions() map[string]getopt.OptionValue {
	optionDefinition := getopt.Options{
		"Beartrap v0.3 by Chris Benedict <chrisbdaemon@gmail.com>",
		getopt.Definitions{
			{"config|c|BEARTRAP_CONFIG", "configuration file", getopt.Required, ""},
		},
	}

	options, _, _, err := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]

	if err != nil || wantsHelp {
		exit_code := 0

		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", err.Error(), "\n", optionDefinition.Help())
			exit_code = err.ErrorCode
		}
		os.Exit(exit_code)
	}

	return options
}
