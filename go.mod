module github.com/fido-device-onboard/go-fdo-client

go 1.23.0

replace github.com/fido-device-onboard/go-fdo/sqlite => ./protocol/sqlite

replace github.com/fido-device-onboard/go-fdo => ./protocol/

replace github.com/fido-device-onboard/go-fdo/fsim => ./protocol/fsim

require (
	github.com/fido-device-onboard/go-fdo v0.0.0-00010101000000-000000000000
	github.com/fido-device-onboard/go-fdo/fsim v0.0.0-00010101000000-000000000000
	hermannm.dev/devlog v0.4.1
)

require (
	github.com/neilotoole/jsoncolor v0.7.1 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/term v0.23.0 // indirect
)
