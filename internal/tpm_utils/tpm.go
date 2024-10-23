// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

//go:build !tpmsim

package tpm_utils

import (
	"fmt"
	"slices"

	"github.com/fido-device-onboard/go-fdo/tpm"
)

var TPMDEVICES = []string{"/dev/tpm0", "/dev/tpmrm0"}

func TpmOpen(tpmPath string) (tpm.Closer, error) {
	if tpmPath == "simulator" {
		return nil, fmt.Errorf("tpm simulator is not supported")
	}

	if slices.Contains(TPMDEVICES, tpmPath) {
		return tpm.Open(tpmPath)
	}

	return nil, fmt.Errorf("invalid tpm path")
}
