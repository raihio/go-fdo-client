// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

package tpm_utils

import (
	"fmt"

	"github.com/fido-device-onboard/go-fdo/tpm"
	"github.com/google/go-tpm/tpm2/transport/linuxtpm"
)

func TpmOpen(tpmPath string) (tpm.Closer, error) {
	if tpmPath == "simulator" {
		return nil, fmt.Errorf("tpm simulator is not supported")
	}

	return linuxtpm.Open(tpmPath)
}
