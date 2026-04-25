# envctl

A CLI tool to manage and sync environment variable sets across local and remote environments using encrypted vaults.

---

## Installation

```bash
go install github.com/yourusername/envctl@latest
```

Or download a pre-built binary from the [Releases](https://github.com/yourusername/envctl/releases) page.

---

## Usage

```bash
# Initialize a new vault in the current directory
envctl init

# Add an environment variable to a named set
envctl set production DATABASE_URL=postgres://user:pass@host/db

# List all variables in a set
envctl list production

# Push a local set to a remote vault
envctl push production

# Pull a remote set and export to your shell
eval $(envctl pull production)

# Remove a variable from a set
envctl unset production DATABASE_URL
```

Environment sets are encrypted at rest using AES-256-GCM. Remote vaults can be backed by S3, GCS, or a self-hosted endpoint configured in `~/.envctl/config.yaml`.

---

## Configuration

```yaml
remote:
  provider: s3
  bucket: my-envctl-vault
  region: us-east-1
encryption:
  key_source: env  # reads ENVCTL_MASTER_KEY
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE) © 2024 Your Name