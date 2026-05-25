### Requirement: List parts in the LEGO catalog

The CLI SHALL provide a command that returns the full paginated list of parts from `GET /api/v3/lego/parts/`. The command MUST follow pagination until all pages are exhausted and MUST emit the aggregated JSON results to stdout. Authentication SHALL use the API key only (no user token), consistent with the existing `lego sets` commands.

#### Scenario: List all parts without filters

- **WHEN** the user runs `cli lego parts get`
- **THEN** the CLI issues authenticated GET requests against `/api/v3/lego/parts/`, following the `next` link until exhausted, and prints the combined `results[]` array as JSON to stdout with a zero exit code

#### Scenario: Filter parts by query parameters

- **WHEN** the user runs `cli lego parts get` with one or more of the supported filter flags (`--part_num`, `--part_nums`, `--part_cat_id`, `--color_id`, `--bricklink_id`, `--brickowl_id`, `--lego_id`, `--ldraw_id`, `--search`, `--ordering`)
- **THEN** the CLI forwards the provided values as URL query parameters on every paginated request, and omits filter parameters that the user did not supply

#### Scenario: Upstream API returns a non-2xx response

- **WHEN** the Rebrickable API responds with a non-2xx status during list pagination
- **THEN** the CLI exits with a non-zero status and writes an error message to stderr identifying the failing URL and status code, and does not emit a partial JSON result to stdout

---

### Requirement: Get a single part by part number

The CLI SHALL provide a command that fetches a single part from `GET /api/v3/lego/parts/{part_num}/` and prints the raw JSON response to stdout.

#### Scenario: Fetch an existing part

- **WHEN** the user runs `cli lego parts getOne --part_num <num>` for a part that exists
- **THEN** the CLI issues `GET /api/v3/lego/parts/<num>/`, writes the JSON response body to stdout, and exits zero

#### Scenario: Fetch a non-existent part

- **WHEN** the user runs `cli lego parts getOne --part_num <num>` and the API responds with 404
- **THEN** the CLI prints a `Part <num> not found` message to stdout and exits zero (consistent with the existing 404 handling for `sets delete` and `setlists delete`)

---

### Requirement: List all colors a part has appeared in

The CLI SHALL provide a command that lists every color a given part has appeared in, from `GET /api/v3/lego/parts/{part_num}/colors/`. The command MUST paginate the response and print the aggregated `results[]` as JSON.

#### Scenario: List colors for a known part

- **WHEN** the user runs `cli lego parts colors --part_num <num>`
- **THEN** the CLI walks pagination on `/api/v3/lego/parts/<num>/colors/` and prints the combined results to stdout

---

### Requirement: Get a single part/color combination

The CLI SHALL provide a command that fetches details for a specific `(part_num, color_id)` pair from `GET /api/v3/lego/parts/{part_num}/colors/{color_id}/`.

#### Scenario: Fetch an existing part/color combination

- **WHEN** the user runs `cli lego parts colorDetail --part_num <num> --color_id <id>` for a valid combination
- **THEN** the CLI fetches the single resource and writes the JSON body to stdout

#### Scenario: Combination does not exist

- **WHEN** the API responds with 404 for the combination
- **THEN** the CLI prints a clear `Part <num> in color <id> not found` message to stdout and exits zero

---

### Requirement: List sets containing a part/color combination

The CLI SHALL provide a command that lists all sets in which a given `(part_num, color_id)` combination appears, from `GET /api/v3/lego/parts/{part_num}/colors/{color_id}/sets/`. The command MUST paginate the response and print the aggregated `results[]` as JSON.

#### Scenario: List sets for a part/color combination

- **WHEN** the user runs `cli lego parts colorSets --part_num <num> --color_id <id>`
- **THEN** the CLI walks pagination on `/api/v3/lego/parts/<num>/colors/<id>/sets/` and prints the combined results to stdout

---

### Requirement: Reuse existing client and pagination infrastructure

The HTTP client methods supporting these commands SHALL live on the existing `api.Client` struct and SHALL reuse the generic `fetchAllPages[T]` helper for every paginated endpoint. No new client constructor or auth mechanism is introduced.

#### Scenario: A new client method is added for parts

- **WHEN** a developer adds a method for any of the five parts endpoints
- **THEN** the method is declared on `*Client` (not a new type), lives in `cli/cmd/api/lego_parts.go`, and any paginated list method delegates to `fetchAllPages[T]` rather than re-implementing the pagination loop

---

### Requirement: Bazel and command registration

The new commands SHALL be reachable from `cli lego parts ...` and SHALL be wired through Cobra's `AddCommand` chain from the existing `legoCmd`. Bazel `BUILD.bazel` files MUST include the new source files in the relevant `go_library` `srcs` so `bazel build //cli` succeeds.

#### Scenario: Running cli without arguments shows parts in help

- **WHEN** the user runs `cli lego --help`
- **THEN** the help output lists `parts` as an available subcommand alongside `sets`

#### Scenario: Bazel builds the CLI with the new files

- **WHEN** a developer runs `bazel build //cli`
- **THEN** the build succeeds and the resulting binary exposes `cli lego parts get|getOne|colors|colorDetail|colorSets`
