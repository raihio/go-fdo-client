// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

// Package main implements client and server modes.
package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

	if err := validateFlags(); err != nil {
		fmt.Fprintf(os.Stderr, "Validation error: %v\n", err)
		os.Exit(1)
	}

	if err := client(); err != nil {
		fmt.Fprintf(os.Stderr, "client error: %v\n", err)
		os.Exit(2)
	}
}

func validateFlags() error {
	if !isValidPath(blobPath) {
		return fmt.Errorf("invalid blob path: %s", blobPath)
	}

	validCipherSuites := []string{
		"A128GCM", "A192GCM", "A256GCM",
		"AES-CCM-64-128-128", "AES-CCM-64-128-256",
		"COSEAES128CBC", "COSEAES128CTR",
		"COSEAES256CBC", "COSEAES256CTR",
	}
	if !contains(validCipherSuites, cipherSuite) {
		return fmt.Errorf("invalid cipher suite: %s", cipherSuite)
	}

	if dlDir != "" && (!isValidPath(dlDir) || !fileExists(dlDir)) {
		return fmt.Errorf("invalid download directory: %s", dlDir)
	}

	if diURL != "" {
		parsedURL, err := url.ParseRequestURI(diURL)
		if err != nil {
			return fmt.Errorf("invalid DI URL: %s", diURL)
		}
		host, port, err := net.SplitHostPort(parsedURL.Host)
		if err != nil {
			return fmt.Errorf("invalid DI URL: %s", diURL)
		}
		if net.ParseIP(host) == nil && !isValidHostname(host) {
			return fmt.Errorf("invalid hostname: %s", host)
		}
		if port != "" && !isValidPort(port) {
			return fmt.Errorf("invalid port: %s", port)
		}
	}

	validDiKeys := []string{"ec256", "ec384", "rsa2048", "rsa3072"}
	if !contains(validDiKeys, diKey) {
		return fmt.Errorf("invalid DI key: %s", diKey)
	}

	validDiKeyEncs := []string{"x509", "x5chain", "cose"}
	if !contains(validDiKeyEncs, diKeyEnc) {
		return fmt.Errorf("invalid DI key encoding: %s", diKeyEnc)
	}

	validKexSuites := []string{"DHKEXid14", "DHKEXid15", "ASYMKEX2048", "ASYMKEX3072", "ECDH256", "ECDH384"}
	if !contains(validKexSuites, kexSuite) {
		return fmt.Errorf("invalid key exchange suite: %s", kexSuite)
	}

	TPMDEVICES := []string{"/dev/tpm0", "/dev/tpmrm0", "simulator"}
	if tpmPath != "" && !slices.Contains(TPMDEVICES, tpmPath) {
		return fmt.Errorf("invalid TPM path: %s", tpmPath)
	}

	for path := range uploads {
		if !isValidPath(path) {
			return fmt.Errorf("invalid upload path: %s", path)
		}

		if !fileExists(path) {
			return fmt.Errorf("file doesn't exist: %s", path)
		}
	}

	if wgetDir != "" && (!isValidPath(wgetDir) || !fileExists(wgetDir)) {
		return fmt.Errorf("invalid wget directory: %s", wgetDir)
	}

	return nil
}

func isValidPath(p string) bool {
	if p == "" {
		return false
	}
	absPath, err := filepath.Abs(p)
	return err == nil && absPath != ""
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isValidHostname(hostname string) bool {
	if len(hostname) > 255 {
		return false
	}
	for _, part := range strings.Split(hostname, ".") {
		if len(part) == 0 || len(part) > 63 {
			return false
		}
		for _, char := range part {
			if !((char >= 'a' && char <= 'z') ||
				(char >= 'A' && char <= 'Z') ||
				(char >= '0' && char <= '9') ||
				char == '-') {
				return false
			}
		}
		if strings.HasPrefix(part, "-") || strings.HasSuffix(part, "-") {
			return false
		}
	}
	return true
}

func isValidPort(port string) bool {
	for _, char := range port {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}
