// SPDX-FileCopyrightText: (C) 2024 Intel Corporation
// SPDX-License-Identifier: Apache 2.0

package tpm_utils

import (
	"fmt"

	"github.com/google/go-tpm/tpm2"
	"github.com/google/go-tpm/tpm2/transport"
)

// TpmNVDefine defines an NV index with the specified size.
func TpmNVDefine(thetpm transport.TPM, nv tpm2.TPMHandle, dataSize uint16, tpmHashAlg tpm2.TPMAlgID) error {
	def := tpm2.NVDefineSpace{
		AuthHandle: tpm2.TPMRHOwner,
		Auth:       tpm2.TPM2BAuth{},
		PublicInfo: tpm2.New2B(
			tpm2.TPMSNVPublic{
				NVIndex: nv,
				NameAlg: tpmHashAlg,
				Attributes: tpm2.TPMANV{
					OwnerWrite:   true,
					OwnerRead:    true,
					AuthWrite:    true,
					AuthRead:     true,
					NoDA:         true,
					ReadSTClear:  true,
					WriteSTClear: true,
					WriteDefine:  true,
				},
				DataSize: dataSize,
			}),
	}
	if _, err := def.Execute(thetpm); err != nil {
		return fmt.Errorf("calling TPM2_NV_DefineSpace: %v", err)
	}
	return nil
}

// TpmNVUnDefine undefines the specified NV index.
func TpmNVUnDefine(thetpm transport.TPM, nv tpm2.TPMHandle, nvName *tpm2.TPM2BName) error {
	undef := tpm2.NVUndefineSpace{
		AuthHandle: tpm2.TPMRHOwner,
		NVIndex: tpm2.NamedHandle{
			Handle: nv,
			Name:   *nvName,
		},
	}
	if _, err := undef.Execute(thetpm); err != nil {
		return fmt.Errorf("calling TPM2_NV_UndefineSpace: %v", err)
	}
	return nil
}

// TpmNVGetSize retrieves the size of the data stored in the specified NV index.
func TpmNVGetSize(thetpm transport.TPM, nv tpm2.TPMHandle) uint16 {
	readPub := tpm2.NVReadPublic{
		NVIndex: nv,
	}
	readPubRsp, err := readPub.Execute(thetpm)
	if err != nil {
		return 0
	}
	nvPublic, err := readPubRsp.NVPublic.Contents()
	if err != nil {
		return 0
	}
	return nvPublic.DataSize
}

// TpmNVRead reads data from the specified NV index.
func TpmNVRead(thetpm transport.TPM, nv tpm2.TPMHandle) ([]byte, error) {
	readPub := tpm2.NVReadPublic{
		NVIndex: nv,
	}
	readPubRsp, err := readPub.Execute(thetpm)
	if err != nil {
		return nil, fmt.Errorf("calling TPM2_NV_ReadPublic: %v", err)
	}
	nvPublic, err := readPubRsp.NVPublic.Contents()
	if err != nil {
		return nil, fmt.Errorf("getting NV public contents: %v", err)
	}
	dataSize := nvPublic.DataSize

	nvName, err := tpm2.NVName(nvPublic)
	if err != nil {
		return nil, fmt.Errorf("calculating name of NV index: %v", err)
	}

	read := tpm2.NVRead{
		AuthHandle: tpm2.TPMRHOwner,
		NVIndex: tpm2.NamedHandle{
			Handle: nv,
			Name:   *nvName,
		},
		Size:   dataSize,
		Offset: 0,
	}
	readRsp, err := read.Execute(thetpm)
	if err != nil {
		return nil, fmt.Errorf("calling TPM2_NV_Read: %v", err)
	}

	return readRsp.Data.Buffer, nil
}

// TpmNVWrite writes data to the specified NV index, deleting existing data if present.
func TpmNVWrite(thetpm transport.TPM, data []byte, nv tpm2.TPMHandle, tpmHashAlg tpm2.TPMAlgID) error {
	// Check if data exists in the NV index
	dataSize := TpmNVGetSize(thetpm, nv)

	// Calculate the NV index name
	nvName, err := tpm2.NVName(&tpm2.TPMSNVPublic{
		NVIndex: nv,
		NameAlg: tpmHashAlg,
		Attributes: tpm2.TPMANV{
			OwnerWrite:  true,
			OwnerRead:   true,
			AuthWrite:   true,
			AuthRead:    true,
			NoDA:        true,
			ReadSTClear: true,
			WriteDefine: true,
		},
		DataSize: uint16(len(data)),
	})
	if err != nil {
		return fmt.Errorf("calculating name of NV index: %v", err)
	}

	// If data exists, undefine the NV index
	if dataSize > 0 {
		if err := TpmNVUnDefine(thetpm, nv, nvName); err != nil {
			return fmt.Errorf("undefining NV index: %v", err)
		}
	}

	// Define the NV index with the new data size
	if err := TpmNVDefine(thetpm, nv, uint16(len(data)), tpmHashAlg); err != nil {
		return fmt.Errorf("defining NV index: %v", err)
	}

	// Write the new data to the NV index
	write := tpm2.NVWrite{
		AuthHandle: tpm2.TPMRHOwner,
		NVIndex: tpm2.NamedHandle{
			Handle: nv,
			Name:   *nvName,
		},
		Data: tpm2.TPM2BMaxNVBuffer{
			Buffer: data,
		},
		Offset: 0,
	}
	if _, err := write.Execute(thetpm); err != nil {
		return fmt.Errorf("calling TPM2_NV_Write: %v", err)
	}
	return nil
}
