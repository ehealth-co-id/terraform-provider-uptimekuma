# Release Checklist

## Pre-Release Steps

### 1. Ensure Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run unit tests
make test

# Run acceptance tests (requires Uptime Kuma running)
docker compose up -d
# Wait for setup, then:
export UPTIMEKUMA_BASE_URL="http://localhost:8000"
export UPTIMEKUMA_USERNAME="admin"
export UPTIMEKUMA_PASSWORD="admin123"
make testacc
```

### 2. Update Documentation

```bash
# Regenerate provider documentation
make generate

# Verify no uncommitted changes
git diff --exit-code
```

### 3. Update Version References

- Update version in any hardcoded references if applicable
- Update `CHANGELOG.md` or release notes

### 4. Create Release Notes

Create `.release_info.md` file with release notes:

```markdown
## What's Changed

- Feature: Description
- Fix: Description
- Docs: Description

**Full Changelog**: https://github.com/ehealth-co-id/terraform-provider-uptimekuma/compare/v0.x.x...v0.y.y
```

### 5. Commit and Push

```bash
git add .
git commit -m "chore: prepare release vX.Y.Z"
git push origin master
```

## Release Steps

### 6. Create and Push Tag

```bash
# Create annotated tag
git tag -a vX.Y.Z -m "Release vX.Y.Z"

# Push tag to trigger release workflow
git push origin vX.Y.Z
```

### 7. Verify Release

1. Check [GitHub Actions](https://github.com/ehealth-co-id/terraform-provider-uptimekuma/actions) for release workflow status
2. Verify release appears on [Releases page](https://github.com/ehealth-co-id/terraform-provider-uptimekuma/releases)
3. Verify artifacts are attached (binaries for all platforms, checksums)

## Post-Release

### 8. Terraform Registry

The Terraform Registry should automatically pick up new releases from GitHub.
Verify at: https://registry.terraform.io/providers/ehealth-co-id/uptimekuma

## Quick Commands Summary

| Command | Description |
|---------|-------------|
| `make fmt` | Format Go code |
| `make lint` | Run golangci-lint |
| `make test` | Run unit tests |
| `make testacc` | Run acceptance tests |
| `make generate` | Generate provider docs |
| `make build` | Build the provider |
| `make install` | Install the provider locally |
| `make` | Run fmt, lint, install, generate |

## Required Secrets (GitHub Actions)

The release workflow requires these secrets configured in GitHub:

- `GPG_PRIVATE_KEY` - GPG private key for signing
- `GPG_PASSPHRASE` - Passphrase for GPG key
- `PAT_TOKEN` - Personal Access Token for creating releases

