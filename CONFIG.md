# Configuration File Reference

This document describes configuration options available for the FDO client. Configuration files use YAML format.

Command line arguments take precedence over configuration file values. If no `--config` option is specified, the client uses defaults where available.

## Configuration File Location

The configuration file can be specified via the `--config` command line parameter, for example:

```bash
# Using YAML configuration file for device initialization:
go-fdo-client device-init --config /etc/fdo-client-config.yaml

# Using YAML configuration file for onboarding:
go-fdo-client onboard --config config.yaml

# Override config file values with CLI flags:
go-fdo-client device-init --config config.yaml --key ec384

# Server URL can be provided via config or CLI argument:
go-fdo-client device-init http://127.0.0.1:8038 --config config.yaml
```

## Configuration Structure

The configuration file uses a hierarchical structure that defines the following sections:

- Global configuration (applies to all commands)
- `device-init` - Device initialization command configuration
- `onboard` - Onboarding command configuration

## Global Configuration

Global configuration options apply to all commands.

| Key | Type | Description | Required |
|-----|------|-------------|----------|
| `debug` | boolean | Print HTTP contents for debugging | No (default: false) |
| `blob` | string | File path of device credential blob | One of `blob` or `tpm` |
| `tpm` | string | Use a TPM at specified path for device credential secrets | One of `blob` or `tpm` |

**Note**: Exactly one of `blob` or `tpm` must be provided.

## Device Initialization Configuration

The device initialization configuration is under the `device-init` section:

| Key | Type | Description | Required |
|-----|------|-------------|----------|
| `server-url` | string | Device initialization server URL (e.g., `http://127.0.0.1:8038`) | Yes (or provide as CLI argument) |
| `key` | string | Key type for device credential. Options: `ec256`, `ec384`, `rsa2048`, `rsa3072` | Yes |
| `key-enc` | string | Public key encoding to use for manufacturer key. Options: `x509`, `x5chain`, `cose` | No (default: `x509`) |
| `device-info` | string | Device information for device credentials. If not specified, will be gathered from the system | No |
| `device-info-mac` | string | Mac address interface (e.g., `eth0`) for device credentials | No |
| `insecure-tls` | boolean | Skip TLS certificate verification | No (default: false) |

**Note**: `device-info` and `device-info-mac` are mutually exclusive.

## Onboarding Configuration

The onboarding configuration is under the `onboard` section:

| Key | Type | Description | Required |
|-----|------|-------------|----------|
| `key` | string | Key type for device credential. Options: `ec256`, `ec384`, `rsa2048`, `rsa3072` | Yes |
| `kex` | string | Key exchange suite. Options: `DHKEXid14`, `DHKEXid15`, `ASYMKEX2048`, `ASYMKEX3072`, `ECDH256`, `ECDH384` | Yes |
| `cipher` | string | Cipher suite for encryption. Options: `A128GCM`, `A192GCM`, `A256GCM`, `AES-CCM-64-128-128`, `AES-CCM-64-128-256`, `COSEAES128CBC`, `COSEAES128CTR`, `COSEAES256CBC`, `COSEAES256CTR` | No (default: `A128GCM`) |
| `download` | string | Directory to download files into (FSIM disabled if empty) | No |
| `upload` | array of strings | List of directories and files to upload from (FSIM disabled if empty) | No |
| `wget-dir` | string | Directory to wget files into (FSIM disabled if empty) | No |
| `echo-commands` | boolean | Echo all commands received to stdout (FSIM disabled if false) | No (default: false) |
| `insecure-tls` | boolean | Skip TLS certificate verification | No (default: false) |
| `allow-credential-reuse` | boolean | Allow credential reuse protocol during onboarding | No (default: false) |
| `max-serviceinfo-size` | integer | Maximum service info size to receive (0-65535) | No (default: 1300) |
| `resale` | boolean | Perform resale | No (default: false) |
| `to2-retry-delay` | string | Delay between failed TO2 attempts when trying multiple Owner URLs from same RV directive. Format: duration string (e.g., `10s`, `1m`) | No (default: `0s` - disabled) |

## Configuration File Examples

### Device Initialization Configuration

```yaml
debug: true
blob: "cred.bin"

device-init:
  server-url: "http://127.0.0.1:8038"
  key: "ec256"
  key-enc: "x509"
  device-info: "my-device"
  insecure-tls: false
```

Usage:
```bash
# Server URL from config:
go-fdo-client device-init --config config.yaml

# Or override server URL via CLI:
go-fdo-client device-init http://192.168.1.100:8038 --config config.yaml
```

### Device Initialization with TPM

```yaml
debug: false
tpm: "/dev/tpmrm0"

device-init:
  server-url: "http://127.0.0.1:8038"
  key: "ec384"
  key-enc: "x509"
  device-info-mac: "eth0"
```

### Onboarding Configuration

```yaml
debug: true
blob: "cred.bin"

onboard:
  key: "ec256"
  kex: "ECDH256"
  cipher: "A128GCM"
  max-serviceinfo-size: 1300
  allow-credential-reuse: false
  to2-retry-delay: "10s"
```

Usage:
```bash
go-fdo-client onboard --config config.yaml
```

### Onboarding with File Upload/Download

```yaml
blob: "cred.bin"

onboard:
  key: "ec256"
  kex: "ECDH256"
  cipher: "A128GCM"
  download: "/var/fdo/downloads"
  upload:
    - "/var/fdo/uploads/file1.txt"
    - "/var/fdo/uploads/dir1"
  wget-dir: "/var/fdo/wget"
  echo-commands: true
```

### Complete Configuration Example

```yaml
debug: false
blob: "cred.bin"

device-init:
  server-url: "https://manufacturing.example.com:443"
  key: "ec384"
  key-enc: "x509"
  device-info: "device-001"

onboard:
  key: "ec384"
  kex: "ECDH384"
  cipher: "A256GCM"
```

## Notes

### Precedence

Configuration values are applied in the following order (highest to lowest precedence):

1. **Command-line arguments** - Values specified via CLI flags or positional arguments
2. **Configuration file** - Values from the YAML config file
3. **Built-in defaults** - Default values defined by the application

Example:
```bash
# config.yaml has: key: "ec256", server-url: "http://localhost:8038"
# CLI overrides with: http://remote:8038 --key ec384
# Result: http://remote:8038 and ec384 are used (CLI takes precedence)
go-fdo-client device-init http://remote:8038 --config config.yaml --key ec384
```

### Device Initialization Server URL

The server URL for device initialization can be provided in two ways:

1. **As a positional CLI argument**:
   ```bash
   go-fdo-client device-init http://127.0.0.1:8038 --config config.yaml
   ```

2. **In the config file**:
   ```yaml
   device-init:
     server-url: "http://127.0.0.1:8038"
   ```

If both are provided, the CLI argument takes precedence.

### Duration Format

The `to2-retry-delay` field accepts Go duration strings:
- `0s` - Disabled (no delay)
- `10s` - 10 seconds
- `1m` - 1 minute
- `1m30s` - 1 minute 30 seconds
- `1h` - 1 hour

### Upload Paths

The `upload` field accepts an array of file and directory paths:
```yaml
upload:
  - "/path/to/file.txt"
  - "/path/to/directory"
  - "/another/path"
```

