---
name: reduce-code
description: >
  Review a codebase or specific files for unnecessary abstractions, packages, features,
  and boilerplate that can be removed to reduce total lines while preserving all required
  functionality. Operates on the principle that code which doesn't need to exist is the
  best code. Produces a prioritized report ordered by reduction impact (largest savings
  first). Use this skill whenever the user says "reduce code", "shrink codebase",
  "find dead code", "simplify design", "remove unnecessary abstraction", "find bloat",
  "how can I cut lines", "コードを削減", "コード行数を減らして", "不要な抽象化を探して",
  "コードを整理して", "コードをシンプルにして".
---

# Reduce Code

Find and remove code that doesn't need to exist, ordered by how much it saves.

The core idea: most bloat comes from premature abstraction — drawing package/module
boundaries, adding flags, or writing utilities before anyone actually needs them. These
decisions don't just add the code they introduce; they force all the surrounding
boilerplate (exports, interfaces, imports, tests for public APIs) to exist too. Removing
the wrong abstraction at the right time saves hundreds or thousands of lines, not just
the abstracted code itself.

Dead code accumulates like weeds: any large codebase that's been around a while has
orphaned lines, functions, and subsystems — code that was called once but no longer is.
**Code that isn't running doesn't work** (Chris Zimmerman). Weeds don't stay small.
Pull them now.

**As Simple as Possible, but No Simpler**: 余分な複雑さは作らない。ただし問題が本質的に
複雑なら単純化しない。削減の判断基準はこの一文に尽きる（Chris Zimmerman, Rule 1）。

**Rule of Three（コード重複）**: 同じロジックの重複が 3 箇所以上になって初めて関数として
抽出する。2 箇所の重複はまだ抽象化のタイミングではない — 3 例目で初めて適切な抽象の形が
見えてくることが多い（Martin Fowler / Don Roberts）。

**YAGNI（インターフェース・抽象型）**: インターフェースは 2 つ以上の具体的な実装が存在するか
直近で必要になるまで定義しない。実装が 1 つしかないインターフェースは将来の拡張を見越した
過剰設計（Speculative Generality）であり、削除して具体型を直接使うべき。

## Workflow

### Step 1: Identify the target and language

- If the user specified a path, use it.
- Otherwise, inspect the repository structure:
  ```
  git ls-files | head -60
  ```
  From the file extensions and directory layout, determine the primary language(s).
  Show the user a summary and confirm the scope before proceeding.

### Step 2: High-impact checks (package/module boundaries)

These are the checks that save hundreds to thousands of lines. Do them first.

**Check: Internal packages or modules with no external consumers**

Find packages/modules/directories that are only imported by code within this same
repository (or, stronger, only within this same binary/app):

```
rg "import.*<package-path>" --type go   # Go
rg "from ['\"]<module>['\"]"             # Python/JS/TS
rg "use <module>"                         # Rust
```

If no outside code imports a package, the package boundary has no enforcement value.
Every class/function/type must be exported (capitalised/public) to cross the boundary,
and every such export needs its own test surface. Removing the boundary eliminates all
of that overhead.

**Before proposing consolidation**, verify:
1. Run the import search above — confirm zero external consumers.
2. Check if the package is referenced in README, API docs, or published as a library.

**If prerequisites pass** → propose consolidation. Note estimated savings as:
package declaration lines + import block lines + re-exported type/interface definitions
+ tests written specifically for the public API.

Also examine the test files for those packages. Tests of a public API may become tests
of private functions after consolidation — they may need rewriting or can be dropped.
Highlight this tradeoff explicitly.

**If prerequisites fail** (e.g., external repo imports this package) → record the
package in the "Verified Unchangeable" output section and move on.

**When to keep a package boundary:**
- Other repositories or binaries import this package
- The package has a stable, versioned public API
- The package isolates a genuinely different concern (e.g., a separate daemon)

### Step 3: Feature/flag removal candidates

Look for CLI flags, subcommands, config options, or code paths that:
- Write data back to a source (write-back, export, sync features)
- Provide a second mode that duplicates the primary mode with tweaks
- Are documented as "future" or "experimental" with no active users

These often carry their own data structures, tests, and error paths.

**Before proposing removal**, verify:
1. Search README/docs for the flag/command by name.
2. Check git log for recent usage: `git log --all --oneline -- <relevant file>`.
3. If the feature appears unused or purely additive, propose removal with estimated savings.

### Step 4: Small-scale mechanical checks

These individually save only a few lines each, but are easy to apply:

**a. Functions with fewer than 3 call sites (Rule of Three)**
A function abstraction earns its keep only when 3 or more independent callers benefit.
If a function is called from 1 or 2 places, inline it unless the name adds essential
clarity that the body cannot convey on its own.

How to find: `rg -c "<function-name>"` — count of 2 (definition + 1 call) or 3
(definition + 2 calls) marks a candidate. Inline both cases unless the name genuinely
encodes a non-obvious invariant.

**b. 実行されていないコード（Code That Isn't Running Doesn't Work）**
呼び出し元がなくなった関数・到達不能な分岐・無効化された機能フラグは、実質壊れている。
削除してもテストが通るなら存在価値がない。

How to find:
- `rg -l "<function-name>"` でファイル数が 1（定義のみ）
- フラグや設定値が常に同じ値に固定されている分岐
- テストからしか呼ばれていない関数（本番コードで未使用）

**c. Utility functions replaceable by the standard library**
Generic helpers that re-implement something the language's standard library already
provides. Common examples are in `references/stdlib-replacements.md` by language.

**d. 実装が 1 つしかないインターフェース（YAGNI / Speculative Generality）**
インターフェースはポリモーフィズムのために存在する。具体的な実装が 1 つしかなければ
ポリモーフィズムはまだ不要 — YAGNI。削除して具体型を直接使う。

探し方: インターフェースのメソッドシグネチャ（または `implements` キーワード）で実装型を検索。
1 つしか見つからなければ削除候補。

**e. Unnecessary type aliases and unused fields**
Type aliases that just rename a primitive, and struct/object fields that are always
zero-valued or always recomputed from other fields.

**f. Language-specific compact patterns**
See `references/go-tactics.md` for Go-specific patterns (single-line error checks, etc.).
Other languages: apply equivalent idioms where the formatter preserves them.

### Step 5: Output the report

Print the report in this format. Sort sections by estimated savings descending.

```
## Reduction Review: <target>

### Summary

| Check | Estimated savings | Risk | Prerequisite |
|---|---|---|---|
| `utils/` module consolidation | ~200 lines | Medium | No external imports (verified) |
| `--export` flag removal | ~60 lines | Low | Confirm unused in docs |
| Single-call functions (4 found) | ~18 lines | Low | None |
| `format_date` → stdlib | ~8 lines | Low | Check Python version ≥ 3.2 |

### High-Impact Proposals

#### Consolidate `utils/` into main module
**Estimated savings**: ~200 lines
**Why this saves so much**: The module boundary forces 12 functions to be exported.
None are imported outside this repository. Collapsing lets everything become
module-private and removes the `__init__.py` re-export surface.
**Prerequisite check**: `rg "from utils import"` — only found in this repo.
**Suggested approach**: Move all files into the main package, prefix private names
with `_`, delete the now-empty `utils/` directory.

...

### Low-Impact Proposals

#### Inline `_format_date(dt)` (line 34 in helpers.py)
Called once at line 118. The inline form `dt.strftime("%Y-%m-%d")` is equally readable.
Saves 6 lines.

...

### Verified Unchangeable
(Items examined where prerequisites confirmed the change is not safe — explain briefly.)

### Not Worth Changing
(Items examined but not worth the risk/effort — explain briefly. Omit this section if empty.)
```

---

## References

- [references/go-tactics.md](references/go-tactics.md) — Go-specific reduction patterns
- [references/stdlib-replacements.md](references/stdlib-replacements.md) — Standard library replacements by language
