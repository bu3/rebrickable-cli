# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Test

Always use Bazel — never `go test` or `go build` directly.

```bash
bazel test //...                              # run all tests
bazel test //cli/cmd/api:api_test             # run unit tests only
bazel test //cli:cli_test                     # run integration tests only
bazel test //... --test_output=all            # show full test output
bazel build //cli                             # build the binary
```

Credentials must be in the environment for integration tests (`//cli:cli_test`). They are loaded automatically from `.envrc` via direnv, and forwarded to Bazel's sandbox via `.bazelrc`.

## Architecture

The CLI is a [Cobra](https://github.com/spf13/cobra) application with three layers:

**`cli/cmd/api/api.go`** — pure HTTP client. All Rebrickable API calls live here as methods on `Client`. Uses [resty](https://github.com/go-resty/resty). `NewClient(apiKey, authToken)` is the public constructor; `newClientWithBaseURL(apiKey, authToken, baseURL)` is used in tests to inject a mock server.

**`cli/cmd/sets.go`** — Cobra command definitions wired to `api.Client` methods. All set and setlist commands live here. `newAPIClient(cmd)` extracts auth from the command context set by the login middleware.

**`cli/cmd/user.go`** — Login middleware via `PersistentPreRunE` on the `user` command. Calls `POST /users/_token/` and stores `user_token` and `api_key` in the command context so all subcommands can access them.

### Auth flow

Every `user *` command triggers `PersistentPreRunE` in `user.go`, which reads `REBRICKABLE_USERNAME`, `REBRICKABLE_PASSWORD`, `REBRICKABLE_API_KEY` from env, calls the login endpoint, and stores the resulting token in `cmd.Context()` under the keys `AuthToken` and `ApiKey`.

### Adding a new API endpoint

1. Add a method to `Client` in `cli/cmd/api/api.go` following the existing pattern.
2. Add a corresponding `*cobra.Command` var in `cli/cmd/sets.go` and register it in the appropriate `*Commands()` init function.
3. Add unit tests in `cli/cmd/api/api_test.go` using `httptest.NewServer` + `newClientWithBaseURL`.
4. Add a txtar integration test in `testdata/` only if the operation is idempotent and self-cleaning (i.e. uses stable IDs like set numbers). Resources with auto-generated IDs (setlists, part lists) are not suitable for txtar tests — cover with unit tests only.

### Remaining API

See `requirements.md` for the full list of unimplemented endpoints from `openapi.spec.json`.
