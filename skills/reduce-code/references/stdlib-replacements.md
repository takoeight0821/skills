# Standard Library Replacements by Language

Common utility functions that reimplemented something already in the standard library.
Check the language version before recommending — not all versions have everything listed.

## Go (1.21+)

| Custom function | Stdlib replacement | Min version |
|---|---|---|
| `sortedKeys(map)` | `slices.Sorted(maps.Keys(m))` | 1.23 |
| `contains(slice, item)` | `slices.Contains(slice, item)` | 1.21 |
| `map(slice, fn)` | no stdlib equivalent — keep | — |
| `filter(slice, fn)` | no stdlib equivalent — keep | — |
| `min(a, b int)` | `min(a, b)` (builtin) | 1.21 |
| `max(a, b int)` | `max(a, b)` (builtin) | 1.21 |
| `clamp(v, lo, hi)` | `max(lo, min(hi, v))` | 1.21 |
| `keys(map)` | `maps.Keys(m)` (returns iter) | 1.23 |
| `values(map)` | `maps.Values(m)` (returns iter) | 1.23 |

Check `go.mod`: `go 1.XX` line must be ≥ the Min version column.

## Python (3.10+)

| Custom function | Stdlib replacement | Min version |
|---|---|---|
| `flatten(list_of_lists)` | `itertools.chain.from_iterable(...)` | all |
| `chunk(lst, n)` | `itertools.batched(lst, n)` | 3.12 |
| `first(iterable)` | `next(iter(iterable))` | all |
| Custom `dataclass` with `__init__` | `@dataclass` | 3.7 |
| Manual `defaultdict` init | `collections.defaultdict` | all |

## TypeScript / JavaScript

| Custom function | Stdlib replacement | Notes |
|---|---|---|
| `groupBy(arr, fn)` | `Object.groupBy(arr, fn)` | ES2024, Node 21+ |
| `flatten(arr)` | `arr.flat()` | ES2019 |
| `unique(arr)` | `[...new Set(arr)]` | all |
| `zip(a, b)` | no stdlib — keep | — |
| `chunk(arr, n)` | no stdlib — keep | — |

## Rust

| Custom function | Stdlib replacement | Notes |
|---|---|---|
| Manual `Option` chain | `.map()`, `.and_then()`, `?` | all |
| `sort_by_key_rev` | `.sort_by(|a, b| b.key.cmp(&a.key))` | all |
| Custom error type with string | `anyhow::Error` / `thiserror` | crate |
