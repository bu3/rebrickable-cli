# Rebrickable Go Client

Client to interact with [Rebrickable API](https://rebrickable.com/api/v3/docs/?key=).

## Machine setup
Set following Environment variables:

```
export REBRICKABLE_USERNAME=...
export REBRICKABLE_PASSWORD="..." #using quotes to prevent issue with values
export REBRICKABLE_API_KEY=....
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
