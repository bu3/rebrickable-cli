# Rebrickable Go Client

Client to interact with [Rebrickable API](https://rebrickable.com/api/v3/docs/?key=).

## Setup

Set the following environment variables:

```
export REBRICKABLE_USERNAME=...
export REBRICKABLE_PASSWORD="..."
export REBRICKABLE_API_KEY=...
```

Build and test with Bazel:

```bash
bazel test //...          # run all tests
bazel build //cli         # build the binary
```

## Supported API

### Sets (`user sets`)

| Command | Flag(s) | API |
|---------|---------|-----|
| `user sets get` | | GET `/users/{token}/sets/` |
| `user sets getOne` | `-n <set_num>` | GET `/users/{token}/sets/{set_num}/` |
| `user sets set` | `-n <set_num>` | POST `/users/{token}/sets/` |
| `user sets replace` | `-n <set_num>` `-q <quantity>` | PUT `/users/{token}/sets/{set_num}/` |
| `user sets delete` | `-n <set_num>` | DELETE `/users/{token}/sets/{set_num}/` |

### Set Lists (`user setLists`)

| Command | Flag(s) | API |
|---------|---------|-----|
| `user setLists get` | | GET `/users/{token}/setlists/` |
| `user setLists getOne` | `-l <list_id>` | GET `/users/{token}/setlists/{list_id}/` |
| `user setLists set` | `-n <name>` | POST `/users/{token}/setlists/` |
| `user setLists update` | `-l <list_id>` `-n <name>` | PATCH `/users/{token}/setlists/{list_id}/` |
| `user setLists replace` | `-l <list_id>` `-n <name>` | PUT `/users/{token}/setlists/{list_id}/` |
| `user setLists delete` | `-l <list_id>` | DELETE `/users/{token}/setlists/{list_id}/` |

### Sets within a Set List (`user setListSets`)

| Command | Flag(s) | API |
|---------|---------|-----|
| `user setListSets get` | `-l <list_id>` | GET `/users/{token}/setlists/{list_id}/sets/` |
| `user setListSets getOne` | `-l <list_id>` `-n <set_num>` | GET `/users/{token}/setlists/{list_id}/sets/{set_num}/` |
| `user setListSets set` | `-l <list_id>` `-n <set_num>` | POST `/users/{token}/setlists/{list_id}/sets/` |
| `user setListSets delete` | `-l <list_id>` `-n <set_num>` | DELETE `/users/{token}/setlists/{list_id}/sets/{set_num}/` |

### LEGO Catalog Sets (`lego sets`)

Requires only `REBRICKABLE_API_KEY` — no login needed.

| Command | Flag(s) | API |
|---------|---------|-----|
| `lego sets list` | | GET `/lego/sets/` |
| `lego sets get` | `-n <set_num>` | GET `/lego/sets/{set_num}/` |
| `lego sets alternates` | `-n <set_num>` | GET `/lego/sets/{set_num}/alternates/` |
| `lego sets minifigs` | `-n <set_num>` | GET `/lego/sets/{set_num}/minifigs/` |
| `lego sets parts` | `-n <set_num>` | GET `/lego/sets/{set_num}/parts/` |
| `lego sets sets` | `-n <set_num>` | GET `/lego/sets/{set_num}/sets/` |

### LEGO Catalog Parts (`lego parts`)

Requires only `REBRICKABLE_API_KEY` — no login needed.

| Command | Flag(s) | API |
|---------|---------|-----|
| `lego parts list` | `--part_num`, `--part_nums`, `--part_cat_id`, `--color_id`, `--bricklink_id`, `--brickowl_id`, `--lego_id`, `--ldraw_id`, `--ordering`, `--search` (all optional) | GET `/lego/parts/` |
| `lego parts get` | `-n <part_num>` | GET `/lego/parts/{part_num}/` |
| `lego parts colors` | `-n <part_num>` | GET `/lego/parts/{part_num}/colors/` |
| `lego parts colorDetail` | `-n <part_num>` `-c <color_id>` | GET `/lego/parts/{part_num}/colors/{color_id}/` |
| `lego parts colorSets` | `-n <part_num>` `-c <color_id>` | GET `/lego/parts/{part_num}/colors/{color_id}/sets/` |
