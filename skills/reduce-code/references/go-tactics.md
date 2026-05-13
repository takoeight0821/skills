# Go-Specific Reduction Tactics

These patterns are specific to Go. Use them in Step 5 of the reduce-code skill.

## 1. `internal/` Package Semantics

Go's `internal/` directory has a special rule: only code in the parent tree can import it.
This is different from "no one imports it" — it's a language-enforced boundary within
the module. Before proposing consolidation:

```bash
# Check who imports it
rg '"<module-path>/internal/<pkg>"'
```

If the only consumers are in the same binary (e.g., all under `cmd/mytool/`), the
internal package has no enforcement value beyond Go's own linker. Collapsing into
`package main` removes:
- `package` declaration
- All export-capitalisation requirements (→ unexport everything)
- Any interfaces that were defined purely to satisfy callers
- Tests written against the exported API

## 2. Single-Line Error Checks (gofmt-compatible)

`gofmt` does not wrap lines by length. This means compact single-line `if` blocks are
preserved exactly as written. Use this for repetitive error checks:

```go
// Before: 4 lines per check
if e.Name == "" {
    return nil, fmt.Errorf("name is required")
}

// After: 1 line per check (gofmt keeps this)
if e.Name == "" { return nil, fmt.Errorf("name is required") }
```

This is appropriate when:
- Multiple sequential guard checks all follow the same pattern
- The body is a single `return` or `continue` statement
- There's no meaningful logic inside the block

Do NOT apply this to blocks with multiple statements or complex logic.

## 3. Standard Library Replacements (Go 1.21+)

### `sortedKeys` / `keys` helpers → `slices` + `maps`

```go
// Before: custom generic function (~7 lines)
func sortedKeys[K ~string, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    slices.Sort(keys)
    return keys
}

// After: inline with stdlib (1 line)
slices.Sorted(maps.Keys(m))
```

Requires: Go 1.23+ for `slices.Sorted` and `maps.Keys` (both became stable in 1.21/1.23).
Check `go.mod` for the `go` directive before recommending.

### `min` / `max` / `clamp`

```go
// Before
func clamp(v, lo, hi int) int {
    if v < lo { return lo }
    if v > hi { return hi }
    return v
}

// After: Go 1.21+ built-ins
min(hi, max(lo, v))
```

### `slices.Contains` replaces manual loop

```go
// Before
for _, item := range list {
    if item == target { return true }
}
return false

// After
slices.Contains(list, target)
```

## 4. Type Aliases That Add No Value

```go
// These add lines but no type safety or clarity:
type reason = string           // untyped alias — just use string
type orgData = map[string]...  // local variable alias — just use the full type

// Keep type aliases only when:
// - They appear in 3+ places and the full type is long/complex
// - They carry semantic meaning enforced elsewhere (e.g., type ID = string in a domain model)
```

## 5. `--write` / Write-Back Feature Patterns

Tools that read data and optionally write it back tend to carry:
- A byte-slice splice or JSON re-serialisation path
- Flag parsing for `--write` / `--dry-run`
- Tests that verify the write-back is reversible

If the tool is used read-only (reporting, auditing), removing the write-back path
eliminates all three. Estimated savings: 50–400 lines depending on how complex the
serialisation logic is.

