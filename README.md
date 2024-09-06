# FIDO Device Onboard - Go Client

`go-fdo-client` is a client implementation of FIDO Device Onboard specification in Go.

[fdo]: https://fidoalliance.org/specs/FDO/FIDO-Device-Onboard-PS-v1.1-20220419/FIDO-Device-Onboard-PS-v1.1-20220419.html
[cbor]: https://www.rfc-editor.org/rfc/rfc8949.html
[cose]: https://datatracker.ietf.org/doc/html/rfc8152

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
  -debug
        Print HTTP contents
  -di URL
        HTTP base URL for DI server
  -download dir
        A dir to download files into (FSIM disabled if empty)
  -insecure-tls
        Skip TLS certificate verification
  -print
        Print device credential blob and stop
  -rv-only
        Perform TO1 then stop
  -upload files
        List of dirs and files to upload files from, comma-separated and/or flag provided multiple times (FSIM disabled if empty)
```


