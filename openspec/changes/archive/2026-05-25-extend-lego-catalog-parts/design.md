## Context

The CLI already follows a stable three-layer pattern for catalog endpoints (see `cli/cmd/api/lego.go`, `cli/cmd/lego_sets.go`, `cli/cmd/lego.go`):

1. **API layer** (`cli/cmd/api/*.go`): methods on `*Client`, one file per resource group (`lego.go` = sets, `user.go` = user resources). Paginated list methods delegate to `fetchAllPages[T]`.
2. **Command layer** (`cli/cmd/*.go`): Cobra `*cobra.Command` vars registered in an `init` function. One file per resource group.
3. **Bazel**: `cli/cmd/BUILD.bazel` and `cli/cmd/api/BUILD.bazel` declare `srcs` explicitly — new files must be added.

Existing infrastructure that this change reuses unchanged:

- `fetchAllPages[T any]` in `cli/cmd/api/api.go` — generic, walks `next` links, returns `(count, []T, error)`.
- `NewLegoClient(apiKey)` — API-key-only constructor (no user token) for catalog endpoints. The parts endpoints have the same auth profile as `lego sets`, so this client is reused as-is.
- `Part` and `PartColor` types in `api.go` — already defined for `SetPart`. The catalog parts endpoints return a richer payload, so we extend with a new `PartDetail` rather than mutating the lean `Part` used by `SetPart`.

## Goals / Non-Goals

**Goals:**

- Implement the 5 `/api/v3/lego/parts/*` endpoints behind 5 `*Client` methods.
- Expose them as Cobra subcommands under `cli lego parts ...` mirroring the verbs used by `cli lego sets ...` (`get`, `getOne`, plus the parts-specific `colors`, `colorDetail`, `colorSets`).
- Preserve the existing 404-as-warning behavior already used by other read commands so scripts and txtar tests can rely on consistent output.
- Add unit-test coverage via `httptest.NewServer` for every new method.
- Add one txtar integration test covering `getOne` for a known stable `part_num`.

**Non-Goals:**

- Implementing the rest of the LEGO catalog (colors, themes, minifigs, part-categories, elements). Those follow as separate changes.
- Adding a `--debug` or `--verbose` flag to print the request URL. The original `Calling URL` log line was already removed from main; restoring it is out of scope.
- Changing the JSON output format (commands continue to emit raw upstream JSON via `fmt.Println(string(output))` like `lego sets get`).
- Caching, retries, rate limiting.

## Decisions

### Decision 1: One new file per layer, named after the resource

- **`cli/cmd/api/lego_parts.go`** — 5 client methods + 2 new types (`PartDetail`, `PartColorDetail`) + their response wrappers.
- **`cli/cmd/lego_parts.go`** — 5 Cobra commands and the `legoPartsCmd` parent, plus an `init()` that wires them under `legoCmd`.

**Alternative considered:** dump everything into the existing `lego.go` / `lego_sets.go`. Rejected — those files would balloon and the file-per-resource convention is already established by `user.go` / `sets.go`.

### Decision 2: Reuse `fetchAllPages[T]` for every list endpoint

All three paginated endpoints (`/parts/`, `/parts/{part_num}/colors/`, `/parts/{part_num}/colors/{color_id}/sets/`) follow the standard `{count, next, previous, results[]}` shape per the OpenAPI spec, so they slot into `fetchAllPages[T]` without modification.

**Alternative considered:** writing bespoke pagination per endpoint to expose `Next` / `Previous` to the user. Rejected — every existing `lego sets` command aggregates all pages and emits a single JSON blob, which is the convention.

### Decision 3: Filter flags on `lego parts get` are independent string flags, not a single `--query` flag

The list endpoint accepts 10 optional query parameters (`part_num`, `part_nums`, `part_cat_id`, `color_id`, `bricklink_id`, `brickowl_id`, `lego_id`, `ldraw_id`, `ordering`, `search`). Each is exposed as its own `StringVarP` Cobra flag, defaulting to empty. The client method takes a `PartsFilter` struct (one field per param) and the Cobra layer constructs it. Empty fields are omitted from the request URL.

**Alternative considered:** a single `--filter key=value` repeatable flag. Rejected — typed flags give better `--help` discoverability and shell completion, matching how the rest of this CLI works.

### Decision 4: `PartDetail` is a new struct, distinct from the existing `Part`

The catalog parts endpoint returns a richer object (`year_from`, `year_to`, `print_of`, `external_ids`, etc.) than the embedded `Part` used inside `SetPart`. Rather than bloating the existing `Part` (and breaking JSON round-trips for set-parts list output), introduce `PartDetail` for the catalog endpoints. Initial fields: those documented in the OpenAPI spec; treat unknown fields as ignored — `resty` + encoding/json drops them by default.

**Alternative considered:** extend `Part` with optional pointers. Rejected — `Part` is widely embedded and changing it has blast radius beyond this change.

### Decision 5: 404 surfaces as a friendly stdout message, not an error

`GetUserSet`, `DeleteUserSet`, etc. already print `Set <num> not found` and return nil on 404. We mirror that for `GetLegoPart` and `GetLegoPartColor` because the catalog is read-only and 404-on-read is a normal user outcome, not a failure. Other status codes still return an error.

**Alternative considered:** return a typed `ErrNotFound`. Rejected for now — would require refactoring existing commands for consistency, and the txtar tests rely on the current stdout-and-exit-zero shape.

### Decision 6: One txtar integration test for `getOne`; the rest are unit-tested

`part_num` is a stable, well-known identifier (e.g. `3001` is a classic 2x4 brick), so `cli lego parts getOne --part_num 3001` is reproducible against the live API. The four other endpoints either depend on parameter combinations (filtered list) or produce volume of output that's tedious to assert against in txtar. Per the convention in `CLAUDE.md`, those get unit tests only.

## Risks / Trade-offs

- **OpenAPI spec is sparse** → response field names are inferred from existing types and the live API. Mitigation: `encoding/json` ignores unknown fields, and the integration test against the live API will catch field-name mismatches on the fields we actually consume.
- **`page_size` defaults are unknown** → `fetchAllPages[T]` doesn't set it, so we get the API default. For a busy endpoint like `/lego/parts/` this could mean many pages and slow output. Mitigation: deferred. If it becomes a problem, add a `--page-size` flag that injects `?page_size=N`.
- **Cobra flag namespace collision** → `--color_id` is used by both `colors get`-style commands and `colorDetail`/`colorSets`. Each command declares its own flags, so no collision occurs.
- **Adding 5 commands at once is more than a minimal slice** → They share infrastructure and types; splitting into 5 PRs would be churn. One PR per resource group is the granularity the repo already uses (see commit `e4b17d8` "add first lego api batch").

## Migration Plan

Purely additive — no rollback or data migration. Steps for the implementer:

1. Add types and `*Client` methods in `cli/cmd/api/lego_parts.go`.
2. Add unit tests in `cli/cmd/api/api_test.go` using `httptest.NewServer` and `newClientWithBaseURL`.
3. Add Cobra commands in `cli/cmd/lego_parts.go` and wire `legoPartsCmd` into `legoCmd` in `cli/cmd/lego.go` or in the new file's `init()`.
4. Update `cli/cmd/api/BUILD.bazel` and `cli/cmd/BUILD.bazel` `srcs`.
5. Add a `testdata/lego_parts.txtar` integration test for `getOne`.
6. Update `requirements.md` to mark these 5 endpoints implemented (35% → 43%).
7. `bazel test //...` must pass.

## Open Questions

None blocking. The exact field set of `PartDetail` will be finalized during implementation by probing the live API with a known part (e.g. `3001`).
