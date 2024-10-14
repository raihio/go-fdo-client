# FIDO Device Onboard - Go Client

`go-fdo-client` is a client implementation of FIDO Device Onboard specification in Go.

[fdo]: https://fidoalliance.org/specs/FDO/FIDO-Device-Onboard-PS-v1.1-20220419/FIDO-Device-Onboard-PS-v1.1-20220419.html
[cbor]: https://www.rfc-editor.org/rfc/rfc8949.html
[cose]: https://datatracker.ietf.org/doc/html/rfc8152

## Prerequisites

- Go 1.23.0 or later
- A Go module initialized with `go mod init`


The `update-deps.sh` script updates all dependencies in your Go module to their latest versions and cleans up the `go.mod` and `go.sum` files.

To update your dependencies, simply run the script:
```sh
./update-deps.sh
```

## Building the Client Application

The client application can be built with `make build` or `go build` directly,

```console
$ make build or go build -o fdo_client ./cmd
$ ./fdo_client

Usage:
  fdo_client [--] [options]

Client options:
  -blob string
        File path of device credential blob (default "cred.bin")
  -cipher suite
        Name of cipher suite to use for encryption (see usage) (default "A128GCM")
  -debug
        Print HTTP contents
  -di URL
        HTTP base URL for DI server
  -di-key string
        Key for device credential [options: ec256, ec384, rsa2048, rsa3072] (default "ec384")
  -di-key-enc string
        Public key encoding to use for manufacturer key [x509,x5chain,cose] (default "x509")
  -download dir
        A dir to download files into (FSIM disabled if empty)
  -echo-commands
        Echo all commands received to stdout (FSIM disabled if false)
  -insecure-tls
        Skip TLS certificate verification
  -kex suite
        Name of cipher suite to use for key exchange (see usage) (default "ECDH384")
  -print
        Print device credential blob and stop
  -rv-only
        Perform TO1 then stop
  -tpm path
        Use a TPM at path for device credential secrets
  -upload files
        List of dirs and files to upload files from, comma-separated and/or flag provided multiple times (FSIM disabled if empty)
  -wget-dir dir
        A dir to wget files into (FSIM disabled if empty)

Key types:
  - RSA2048RESTR
  - RSAPKCS
  - RSAPSS
  - SECP256R1
  - SECP384R1

Encryption suites:
  - A128GCM
  - A192GCM
  - A256GCM
  - AES-CCM-64-128-128 (not implemented)
  - AES-CCM-64-128-256 (not implemented)
  - COSEAES128CBC
  - COSEAES128CTR
  - COSEAES256CBC
  - COSEAES256CTR

Key exchange suites:
  - DHKEXid14
  - DHKEXid15
  - ASYMKEX2048
  - ASYMKEX3072
  - ECDH256
  - ECDH384
```

## Building the Client Application with TPM Simulator
To build the client application with the tpmsim tag, you can use make build-tpmsim or go build wit-tags option:
```console
$ make build-tpmsim
# or
$ go build -tags tpmsim -o fdo_client ./cmd/fdo_client/
$ ./fdo_client -tpm simulator

The TPM simulator may be used with 3 caveats:

1. RSA3072 keys are not supported
2. OpenSSL libraries and headers must be installed
3. The `tpmsim` build tag must be used
```