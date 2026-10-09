<div align="center">

# HOSTAGE LVX

A fast, scope-aware dangling-DNS and takeover engine for real operators.

[![CI](https://img.shields.io/badge/CI-live-cyan?style=for-the-badge&logo=githubactions&logoColor=black)](./.github/workflows/ci.yml)
[![Release](https://img.shields.io/badge/Release-v0.3.0-00d9ff?style=for-the-badge&logo=github&logoColor=black)](../../releases/latest)
[![License](https://img.shields.io/badge/License-MIT-111111?style=for-the-badge)](./LICENSE)
[![Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/leviathan-offsec/HostageLVX?quickstart=1)

</div>

`hostage` resolves a list of targets, fingerprints the HTTP surface, and highlights takeover-prone endpoints with low false-positive noise. It is designed for reusable recon workflows, operator tooling, and CI-friendly scanning.

---

## What it does

- Checks dangling DNS and takeover-prone subdomains
- Fingerprints common takeover services (GitHub Pages, AWS S3, Azure, Heroku, Fastly, CloudFront, Vercel, etc.)
- Uses wildcard-DNS canaries to reduce false positives
- Supports JSONL output for automation and SIEM pipelines
- Runs on Linux, macOS, and Windows

---

## Install

### Via Homebrew (macOS & Linux)

```bash
brew install leviathan-offsec/tap/hostage
hostage --help
```

### One-line installer (curl | sh)


```bash
curl -sSL https://raw.githubusercontent.com/leviathan-offsec/HostageLVX/main/install.sh | sh
```

Install into a specific directory:

```bash
curl -sSL https://raw.githubusercontent.com/leviathan-offsec/HostageLVX/main/install.sh | sh -s -- -b /usr/local/bin
```

### Prebuilt binary

Download from the [Releases](../../releases) page:

- Windows x64
- Linux x64
- Linux arm64
- macOS amd64
- macOS arm64

### Go install

```bash
go install github.com/leviathan-offsec/HostageLVX@latest
```

Requires Go 1.22+.

### Docker

```bash
docker pull ghcr.io/cyeezy08/hostagelvx:latest
```

---

## Quick usage

```bash
# basic scan from a file
hostage @subs.txt

# stdin input
cat subs.txt | hostage -

# JSONL output
hostage -json -o findings.jsonl @subs.txt

# silent mode: only print TAKEOVER / LIKELY results
hostage -silent @subs.txt

# use a custom resolver
hostage -dns-server 1.1.1.1:53 @subs.txt

# print the fingerprint database
hostage -fingerprints
```

### Common flags

| Flag | Default | Description |
|---|---|---|
| `-t` | `50` | Concurrent worker count |
| `-timeout` | `8.0` | DNS + HTTP timeout (seconds) |
| `-resolver` | `doh` | `doh` or `system` |
| `-dns-server` | empty | Custom resolver, e.g. `1.1.1.1:53` |
| `-json` | `false` | Emit JSONL |
| `-silent` | `false` | Only print confirmed takeovers |
| `-o` | stdout | Output file |
| `-fingerprints` | `false` | Show all fingerprints |
| `-no-wildcard-check` | `false` | Skip wildcard detection |
| `-V` | — | Print version |

---

## Output examples

### Text mode

```text
[X] TAKEOVER   legacy.example.com
        CNAME   legacy-example-com.herokuapp.com
        service Heroku
        note    claim by creating a Heroku app with the matching name

[+] ALIVE       api.example.com -> 104.26.14.225
        server  cloudflare
        status  200 OK
        tech    Cloudflare, nginx
```

### JSONL mode

```json
{"host":"legacy.example.com","verdict":"takeover","cname":"legacy-example-com.herokuapp.com","service":"Heroku","status":404,"note":"claim by creating a Heroku app with the matching name"}
{"host":"api.example.com","verdict":"alive","ip":"104.26.14.225","server":"cloudflare","status":200,"tech":["Cloudflare","nginx"]}
```

### Exit codes

- `0` — no takeover candidates found
- `1` — at least one takeover candidate found
- `2` — usage or I/O error

---

## Fingerprints

`hostage` includes a built-in fingerprint set for common services including:

- GitHub Pages
- AWS S3
- Azure Web Apps
- Heroku
- Fastly
- Shopify
- Tumblr
- Zendesk
- Bitbucket Cloud
- Surge.sh
- Readme.io
- Pantheon
- Ghost
- Helpjuice
- CloudFront
- Netlify
- Vercel
- Google Cloud Storage
- DigitalOcean Spaces
- Firebase Hosting
- Webflow

Use `hostage -fingerprints` to print the full database.

---

## Wildcard DNS detection

Hostage LVX can detect wildcard DNS zones and suppress parking-page noise before it becomes a false positive. This is useful when many subdomains resolve to the same placeholder response.

---

## Developer environment

```bash
make build       # build hostage binary
make build-all   # cross-compile all platforms
test
go test ./...   # run tests
make vet         # go vet
make lint        # golangci-lint
make docker      # build docker image
```

### Codespaces

[![Open in Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/leviathan-offsec/HostageLVX?quickstart=1)

---

## Legal / usage

If you cannot produce written authorization for a target, do not scan it.

See [DISCLAIMER.md](./DISCLAIMER.md) for the full policy text.

---

## Contributing

PRs are welcome for new fingerprints, bug fixes, and platform improvements. Keep dependencies minimal and avoid submitting real target data.

---

## License

MIT — see [LICENSE](./LICENSE).

---

Built by [@cyeezy08](https://github.com/cyeezy08) and published under the Leviathan OffSec tooling stack.
