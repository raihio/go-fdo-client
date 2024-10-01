module github.com/fido-device-onboard/go-fdo-client

go 1.23.0

replace github.com/fido-device-onboard/go-fdo/sqlite => ./protocol/sqlite

replace github.com/fido-device-onboard/go-fdo => ./protocol/

replace github.com/fido-device-onboard/go-fdo/fsim => ./protocol/fsim

replace github.com/fido-device-onboard/go-fdo/tpm => ./protocol/tpm

require (
	github.com/fido-device-onboard/go-fdo v0.0.0-00010101000000-000000000000
	github.com/fido-device-onboard/go-fdo/fsim v0.0.0-00010101000000-000000000000
	github.com/fido-device-onboard/go-fdo/tpm v0.0.0-00010101000000-000000000000
	github.com/google/go-tpm v0.9.2-0.20240920144513-364d5f2f78b9
	github.com/google/go-tpm-tools v0.4.4
	hermannm.dev/devlog v0.4.1
)

require (
	github.com/google/go-configfs-tsm v0.3.2 // indirect
	github.com/neilotoole/jsoncolor v0.7.1 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
)
