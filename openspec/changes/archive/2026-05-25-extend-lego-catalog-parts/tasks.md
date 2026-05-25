## 1. API client layer

- [x] 1.1 Create `cli/cmd/api/lego_parts.go` with package declaration and required imports
- [x] 1.2 Add `PartDetail` struct in `api.go` (or `lego_parts.go`) covering the fields returned by `/lego/parts/{part_num}/` (`part_num`, `name`, `part_cat_id`, `part_url`, `part_img_url`, `external_ids`, `print_of`, `year_from`, `year_to`)
- [x] 1.3 Add `PartColorDetail` struct covering the fields returned by `/lego/parts/{part_num}/colors/{color_id}/` (`color_id`, `color_name`, `num_sets`, `num_set_parts`, `part_img_url`, `elements`)
- [x] 1.4 Add `LegoPartsResponse`, `PartColorsResponse` wrapper types following the `{Count, Next, Previous, Results[]}` pattern used by existing `*Response` types
- [x] 1.5 Add `PartsFilter` struct with one string field per supported query param (`PartNum`, `PartNums`, `PartCatID`, `ColorID`, `BricklinkID`, `BrickowlID`, `LegoID`, `LdrawID`, `Ordering`, `Search`)
- [x] 1.6 Implement `(c *Client) GetLegoParts(filter PartsFilter) (*LegoPartsResponse, error)` — uses `fetchAllPages[PartDetail]`, appends non-empty filter fields as query parameters to the first URL
- [x] 1.7 Implement `(c *Client) GetLegoPart(partNum string) (*PartDetail, error)` — single GET; on 404 print `Part <num> not found` and return `(nil, nil)`
- [x] 1.8 Implement `(c *Client) GetLegoPartColors(partNum string) (*PartColorsResponse, error)` — uses `fetchAllPages[PartColorDetail]`
- [x] 1.9 Implement `(c *Client) GetLegoPartColor(partNum, colorID string) (*PartColorDetail, error)` — single GET; on 404 print `Part <num> in color <id> not found` and return `(nil, nil)`
- [x] 1.10 Implement `(c *Client) GetLegoPartColorSets(partNum, colorID string) (*LegoSetsResponse, error)` — uses `fetchAllPages[Set]` and returns the existing `LegoSetsResponse` wrapper

## 2. Cobra command layer

- [x] 2.1 Create `cli/cmd/lego_parts.go` with package, imports, and the `legoPartsCmd` parent (`Use: "parts"`, short description)
- [x] 2.2 Define `getLegoPartsCmd` (`Use: "list"`) with flags for every `PartsFilter` field; in `RunE` build the filter, call `NewLegoClient(apiKey).GetLegoParts(...)`, marshal the response to JSON and `fmt.Println` it _(renamed `get` → `list` to match existing `lego sets list` convention)_
- [x] 2.3 Define `getLegoPartCmd` (`Use: "get"`) with `--part_num` flag; mark it required; print the JSON response (or do nothing extra on 404 — the API method already printed the friendly message) _(renamed `getOne` → `get` to match existing `lego sets get` convention)_
- [x] 2.4 Define `getLegoPartColorsCmd` (`Use: "colors"`) with required `--part_num` flag
- [x] 2.5 Define `getLegoPartColorCmd` (`Use: "colorDetail"`) with required `--part_num` and `--color_id` flags
- [x] 2.6 Define `getLegoPartColorSetsCmd` (`Use: "colorSets"`) with required `--part_num` and `--color_id` flags
- [x] 2.7 In the file's `init()`, register all 5 subcommands on `legoPartsCmd` and register `legoPartsCmd` on `legoCmd`

## 3. Bazel wiring

- [x] 3.1 Add `lego_parts.go` to `srcs` in `cli/cmd/api/BUILD.bazel`
- [x] 3.2 Add `lego_parts.go` to `srcs` in `cli/cmd/BUILD.bazel`
- [x] 3.3 Run `bazel build //cli` and confirm it builds cleanly

## 4. Unit tests

- [x] 4.1 In `cli/cmd/api/api_test.go`, add a table-driven test for `GetLegoParts` covering: empty filter (200), filter with multiple params (verify URL query string), 500 response (error), pagination across 2+ pages
- [x] 4.2 Add a table-driven test for `GetLegoPart` covering: 200 (returns struct), 404 (returns nil + prints message), 500 (error)
- [x] 4.3 Add a table-driven test for `GetLegoPartColors` covering: 200 with pagination, 500
- [x] 4.4 Add a table-driven test for `GetLegoPartColor` covering: 200, 404, 500
- [x] 4.5 Add a table-driven test for `GetLegoPartColorSets` covering: 200 with pagination, 500
- [x] 4.6 Run `bazel test //cli/cmd/api:api_test --test_output=errors` and confirm green

## 5. Integration test

- [x] 5.1 Create `testdata/lego_parts.txtar` with one scenario: `cli lego parts get --part_num 3001` and assert the response contains `"part_num": "3001"` _(used `get` not `getOne` per renamed verb)_
- [x] 5.2 Run `bazel test //cli:cli_test --test_output=errors` and confirm green

## 6. Docs

- [x] 6.1 Update `requirements.md`: under "Parts", change the table header to mark it ✅ implemented and list these endpoints in the same shape used for "Sets"
- [x] 6.2 Update the overall counter at the top of `requirements.md` from "21 of 60" to "26 of 60" (43%)
- [x] 6.3 Update `README.md` if there is a usage example for `lego sets` — mirror it for `lego parts`

## 7. Final verification

- [x] 7.1 Run `bazel test //...` and confirm all targets pass
- [x] 7.2 Run `cli lego parts --help` and confirm all 5 subcommands appear with their flags
- [x] 7.3 Spot-check one live request against the Rebrickable API: `cli lego parts get --part_num 3001`
