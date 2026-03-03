# Rebrickable API — Remaining Endpoints

Endpoints from `openapi.spec.json` not yet implemented. Grouped by resource.

---

## Badges

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/badges/` | List all badges |
| GET | `/api/v3/users/badges/{id}/` | Get a specific badge |

No auth token required — uses API key only.

---

## User Profile

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/profile/` | Get user profile details |

---

## All Parts

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/allparts/` | Get all parts across all sets and part lists |

---

## Parts

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/parts/` | Get parts owned by the user |

---

## Minifigs

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/minifigs/` | Get minifigs owned by the user |

---

## Build

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/build/{set_num}/` | Check if user can build a set from owned parts |

---

## Lost Parts

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/lost_parts/` | List lost parts |
| POST | `/api/v3/users/{user_token}/lost_parts/` | Report a lost part. Body: `inv_part_id` (integer) |
| DELETE | `/api/v3/users/{user_token}/lost_parts/{id}/` | Remove a lost part record |

---

## Part Lists

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/partlists/` | List all part lists |
| POST | `/api/v3/users/{user_token}/partlists/` | Create a part list. Body: `name` (required) |
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/` | Get a specific part list |
| PATCH | `/api/v3/users/{user_token}/partlists/{list_id}/` | Partially update a part list. Body: `name` |
| PUT | `/api/v3/users/{user_token}/partlists/{list_id}/` | Replace a part list. Body: `name` (required) |
| DELETE | `/api/v3/users/{user_token}/partlists/{list_id}/` | Delete a part list |

### Parts within a Part List

| Method | Path | Notes |
|--------|------|-------|
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/parts/` | List parts in a part list |
| POST | `/api/v3/users/{user_token}/partlists/{list_id}/parts/` | Add a part. Body: `part_num`, `color_id`, `quantity` |
| GET | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | Get a specific part |
| PUT | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | Replace a part. Body: `quantity` |
| DELETE | `/api/v3/users/{user_token}/partlists/{list_id}/parts/{part_num}/{color_id}/` | Remove a part |

---

## Sets — Remaining Operations

| Method | Path | Notes |
|--------|------|-------|
| POST | `/api/v3/users/{user_token}/sets/sync/` | Sync sets. Replaces all sets in the default list |

---

## Set Lists — Remaining Operations

| Method | Path | Notes |
|--------|------|-------|
| PATCH | `/api/v3/users/{user_token}/setlists/{list_id}/sets/{set_num}/` | Partially update a set in a list. Body: `quantity`, `include_spares` |
| PUT | `/api/v3/users/{user_token}/setlists/{list_id}/sets/{set_num}/` | Replace a set in a list. Body: `quantity`, `include_spares` |

---

## Implementation Notes

- All user endpoints require both `Authorization: key {api_key}` header and `{user_token}` in the path, obtained via `POST /users/_token/`.
- Paginated list responses follow the `{ count, next, previous, results[] }` shape already used by `SetsResponse` and `SetListsResponse`.
- Part list parts are identified by a composite key `(part_num, color_id)` rather than a single ID.
- `sets/sync/` is a destructive operation — it replaces the entire default set list.
- Integration testing of resources with auto-generated IDs (part lists, set lists) is impractical with txtar; cover with unit tests using httptest mocks instead.
