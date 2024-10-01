// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

// Package main implements client and server modes.
package main

import (
	"flag"
	"fmt"
	"os"
)

var flags = flag.NewFlagSet("root", flag.ContinueOnError)

func main() {
	if err := flags.Parse(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var args []string
	if flags.NArg() > 1 {
		args = flags.Args()[1:]
		if flags.Arg(1) == "--" {
			args = flags.Args()[2:]
		}
	}
	if err := clientFlags.Parse(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := client(); err != nil {
		fmt.Fprintf(os.Stderr, "client error: %v\n", err)
		os.Exit(2)
	}
}
