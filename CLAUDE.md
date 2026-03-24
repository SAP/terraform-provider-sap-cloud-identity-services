# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Terraform provider for SAP Cloud Identity Services (SCI). Built with the **Terraform Plugin Framework** (not the older SDK). Registry address: `sap/sap-cloud-identity-services`. Provider type prefix: `sci_`.

## Common Commands

```bash
make build          # Compile the provider
make install        # Build and install locally (go install)
make test           # Unit tests (uses VCR cassettes, no live API needed)
make testacc        # Acceptance tests (requires TF_ACC=1, set automatically)
make lint           # Run golangci-lint
make fmt            # Format code with gofmt
make generate       # Regenerate documentation via tfplugindocs
```

**Run a single test:**
```bash
go test -v -run TestResourceUser_Create -tags=all -timeout=900s ./sci/provider/
```

**Re-record VCR cassettes** (requires live SCI tenant credentials):
```bash
TEST_RECORD=true SCI_USERNAME=<user> SCI_PASSWORD=<pass> go test -v -run TestResourceUser_Create -tags=all -timeout=900s ./sci/provider/
```

## Architecture

### Package Layout

- `main.go` — Entry point; runs `providerserver.Serve` with protocol v6. Contains `//go:generate` directives for tfplugindocs.
- `sci/provider/` — All provider logic: resources, data sources, type definitions, and tests.
- `internal/cli/` — HTTP client layer. `Client` handles raw HTTP; `SciClient` (in `iasclient.go`) composes resource-specific CLI managers (ApplicationsCli, UsersCli, GroupsCli, SchemasCli, CorporateIdPsCli).
- `internal/cli/apiObjects/` — API request/response structs organized by resource type.
- `internal/utils/` — Custom validators (certificate, datetime, JSON, IDP type) and PATCH request builders.

### Resource Structure Pattern

Each resource follows this file convention in `sci/provider/`:
- `resource_<name>.go` — Resource CRUD implementation
- `datasource_<name>.go` — Singular data source (read one)
- `datasource_<name>s.go` — Plural data source (list all)
- `type_<name>.go` — Terraform schema type definitions (framework types ↔ Go structs)
- `*_test.go` — Tests using VCR cassettes

**Resources:** `sci_application`, `sci_user`, `sci_group`, `sci_schema`, `sci_corporate_idp`

### API Communication

- **SCIM 2.0** (`application/scim+json`) for users and groups — uses PATCH operations
- **REST/JSON** (`application/json`) for applications, schemas, corporate IdPs
- Updates use PATCH (recently migrated from PUT)
- The `Client.Execute()` method handles request building, error parsing, and response deserialization

### Testing Infrastructure

Tests use **go-vcr** to record/replay HTTP interactions. Cassettes are YAML files in `sci/provider/fixtures/`.

Key test helpers in `provider_test.go`:
- `setupVCR(t, cassetteName)` — Initializes recorder; uses `ModeRecordOnce` by default, `ModeRecordOnly` when `TEST_RECORD=true`
- `getTestProviders(httpClient)` — Returns `ProtoV6ProviderFactories`; pass `rec.GetDefaultClient()` for VCR tests, `nil` for error-path unit tests
- `providerConfig(_, testUser)` — Generates HCL provider block with test credentials

Tests with `getTestProviders(nil)` are pure validation tests (no API calls). Tests with `rec.GetDefaultClient()` replay recorded HTTP interactions.

### Authentication

Three mutually exclusive methods, resolved in this priority order in `provider.go`:
1. **OAuth2 Client Credentials** — `client_id` + `client_secret` (or `SCI_CLIENT_ID`/`SCI_CLIENT_SECRET` env vars)
2. **X.509 Certificate** — `p12_certificate_content` + `p12_certificate_password` (or `SCI_P12_CERTIFICATE_PASSWORD` env var)
3. **Basic Auth** — `username` + `password` (or `SCI_USERNAME`/`SCI_PASSWORD` env vars)

## Development Conventions

- **Conventional Commits** required for PR titles: `fix:`, `feat:`, `refactor!:`, `feat(sci_application):`, etc.
- Run `make test` and `make generate` before submitting PRs.
- Documentation in `docs/` is auto-generated — edit templates in `templates/` or resource descriptions in code instead.
- Terraform dev override for local testing: configure `~/.terraformrc` with `dev_overrides` pointing to your `GOBIN`.
