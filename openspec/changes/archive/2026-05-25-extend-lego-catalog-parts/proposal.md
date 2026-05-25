## Why

The CLI currently exposes only the `/lego/sets/*` slice of the Rebrickable LEGO catalog (6 of 22 catalog endpoints). Parts is the next most useful slice for users browsing the catalog: it answers "what is this brick?", "what colors does it come in?", and "which sets contain it?". Adding parts unblocks future work that depends on part identifiers (part lists, build-from-owned-parts, lost parts).

## What Changes

- Add 5 read-only HTTP client methods covering the `/api/v3/lego/parts/*` endpoint family:
  - `GET /api/v3/lego/parts/` (list, with filters: `part_num`, `part_nums`, `part_cat_id`, `color_id`, `bricklink_id`, `brickowl_id`, `lego_id`, `ldraw_id`, `search`, `ordering`)
  - `GET /api/v3/lego/parts/{part_num}/`
  - `GET /api/v3/lego/parts/{part_num}/colors/`
  - `GET /api/v3/lego/parts/{part_num}/colors/{color_id}/`
  - `GET /api/v3/lego/parts/{part_num}/colors/{color_id}/sets/`
- Add a new Cobra subcommand tree under `lego parts` mirroring the existing `lego sets` shape (`get`, `getOne`, `colors`, `colorDetail`, `colorSets`).
- Reuse the existing generic pagination helper (`fetchAllPages[T]`) for all list endpoints.
- Update `requirements.md` to mark these 5 endpoints as implemented.

No breaking changes — purely additive. No auth changes (these are unauthenticated, API-key-only endpoints, same as `lego sets`).

## Capabilities

### New Capabilities
- `lego-catalog-parts`: Read-only access to LEGO parts in the catalog — list parts with rich filters, fetch a specific part, list colors a part appears in, fetch a specific part/color combination, and list sets that contain a given part/color combination.

### Modified Capabilities
<!-- None — `lego-catalog-sets` is not formally specified yet and is unaffected by this change. -->

## Impact

- **New files**: `cli/cmd/api/lego_parts.go` (client methods + types), `cli/cmd/lego_parts.go` (Cobra commands), test additions in `cli/cmd/api/api_test.go`.
- **Modified files**: `cli/cmd/lego.go` (register `parts` subcommand), `requirements.md` (mark endpoints implemented), `cli/cmd/api/BUILD.bazel` and `cli/cmd/BUILD.bazel` (Bazel `srcs`).
- **No DB / external service changes.** No new dependencies.
- **Integration tests**: parts are identified by stable `part_num`, so txtar coverage is appropriate for at least `getOne` and one filter case — following the existing convention for `lego sets`.
