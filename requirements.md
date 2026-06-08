# Rebrickable API — Remaining Endpoints

Derived from `openapi.spec.json`. **26 of 60 operations implemented** (43%).

---

## LEGO Catalog

These are read-only, unauthenticated endpoints (API key only, no user token).

**Implemented:** Parts (5), Sets (6), Colors (2), Elements (1), Minifigs (4), Part Categories (2), Themes (2) — 20 of 20 endpoints.

### Implemented ✅

#### Colors ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/colors/` |
| GET | `/api/v3/lego/colors/{id}/` |

#### Elements ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/elements/{element_id}/` |

#### Minifigs ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/minifigs/` |
| GET | `/api/v3/lego/minifigs/{set_num}/` |
| GET | `/api/v3/lego/minifigs/{set_num}/parts/` |
| GET | `/api/v3/lego/minifigs/{set_num}/sets/` |

#### Part Categories ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/part_categories/` |
| GET | `/api/v3/lego/part_categories/{id}/` |

#### Themes ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/themes/` |
| GET | `/api/v3/lego/themes/{id}/` |

#### Parts ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/parts/` |
| GET | `/api/v3/lego/parts/{part_num}/` |
| GET | `/api/v3/lego/parts/{part_num}/colors/` |
| GET | `/api/v3/lego/parts/{part_num}/colors/{color_id}/` |
| GET | `/api/v3/lego/parts/{part_num}/colors/{color_id}/sets/` |

#### Sets ✅

| Method | Path |
|--------|------|
| GET | `/api/v3/lego/sets/` |
| GET | `/api/v3/lego/sets/{set_num}/` |
| GET | `/api/v3/lego/sets/{set_num}/alternates/` |
| GET | `/api/v3/lego/sets/{set_num}/minifigs/` |
| GET | `/api/v3/lego/sets/{set_num}/parts/` |
| GET | `/api/v3/lego/sets/{set_num}/sets/` |

---

## User Endpoints (not yet implemented)

All require `Authorization: key {api_key}` header and `{user_token}` in the path.
Auth token is obtained via `POST /users/_token/` (already implemented in `user.go`).

### Badges

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/badges/` | No user token required |
| GET | `/api/v3/users/badges/{id}/` | No user token required |

### Profile

| Method | Path |
|--------|------|
| GET | `/api/v3/users/{user_token}/profile/` |

### All Parts

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/allparts/` | Aggregates parts across all sets and part lists |

### Parts

| Method | Path |
|--------|------|
| GET | `/api/v3/users/{user_token}/parts/` |

### Minifigs

| Method | Path |
|--------|------|
| GET | `/api/v3/users/{user_token}/minifigs/` |

### Build

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/build/{set_num}/` | Check if user can build a set from owned parts |

### Lost Parts

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/lost_parts/` | |
| POST | `/api/v3/users/{user_token}/lost_parts/` | Body: `inv_part_id` (integer) |
| DELETE | `/api/v3/users/{user_token}/lost_parts/{id}/` | |

### Part Lists

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/partlists/` | |
| POST | `/api/v3/users/{user_token}/partlists/` | Body: `name` (required) |
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/` | |
| PATCH | `/api/v3/users/{user_token}/partlists/{list_id}/` | Body: `name` |
| PUT | `/api/v3/users/{user_token}/partlists/{list_id}/` | Body: `name` (required) |
| DELETE | `/api/v3/users/{user_token}/partlists/{list_id}/` | |
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/parts/` | |
| POST | `/api/v3/users/{user_token}/partlists/{list_id}/parts/` | Body: `part_num`, `color_id`, `quantity` |
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | |
| PUT | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | Body: `quantity` |
| DELETE | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | |

### Sets — remaining operations

| Method | Path | Notes |
|--------|------|-------|
| POST | `/api/v3/users/{user_token}/sets/sync/` | Destructive — replaces the entire default set list |

### Set List Sets — remaining operations

| Method | Path | Notes |
|--------|------|-------|
| PATCH | `/api/v3/users/{user_token}/setlists/{list_id}/sets/{set_num}/` | Body: `quantity`, `include_spares` |
| PUT | `/api/v3/users/{user_token}/setlists/{list_id}/sets/{set_num}/` | Body: `quantity`, `include_spares` |

---

## Implementation Notes

- LEGO catalog endpoints do not require a user token — they can use a simpler client (API key only, no login step).
- Paginated list responses follow `{ count, next, previous, results[] }` — reuse existing `SetsResponse` / `SetListsResponse` as a pattern.
- Part list parts use a composite key `(part_num, color_id)` rather than a single ID.
- `sets/sync/` is destructive — it replaces the entire default set list.
- Integration testing with txtar is practical only for resources with stable, known identifiers (e.g. set numbers). Resources with auto-generated IDs (part lists, setlists) should be covered by unit tests using `httptest` mocks.
