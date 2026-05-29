# Process-Scaffolding Detection Patterns and Fix Guide

When the logic a document needed *while it was being built* — the instructions it followed, the options it weighed, the consistency checks it ran, the drafting-stage labels, the rhetorical stances it announced — survives into the body text, the reader is forced to read *how the document was made* instead of its subject. You don't leave the scaffolding on a finished building. The same goes for a document.

This guide detects "process scaffolding" and shows how to remove it without losing real content.

## Why This Is a Problem

- **Stance announcements** ("To be honest", "Frankly", "Let me emphasize here") tell the reader about the author's posture, not about the subject. Delete them and the claim is exactly as true.
- **Verification / consistency notes** ("to avoid circularity", "to keep this consistent", "to remove duplication") are the author keeping their own argument's books. The reader needs the conclusion, not the internal-validity ledger.
- **Drafting-stage labels** ("draft", "tentative", "WIP", "たたき台", "(仮)") describe where the document is in its lifecycle, not what it says.

These leave the document readable but make it *about itself*. They differ from temporal expressions (which become stale): scaffolding is wrong the moment it is written, regardless of time.

## The Core Discrimination Test (the load-bearing part)

The regex below has a high false-positive rate. Every hit must pass one question:

> **Is this explaining the subject, or is it explaining how this document was made?**

If deleting the span leaves the claim equally true, it is scaffolding. For the hard boundary cases, split with this rule:

> **A caveat that bears on the reader's actual situation or decision is legitimate content — keep it. The author worrying about their own argument's internal validity (circularity, consistency, duplication) is scaffolding — drop it.**

The decisive question is **whose concern is being served**:

- **Legitimate caveat** — bears on the reader's real situation. e.g. "We tentatively assume the current Read value (final value decided separately)" → the reader must read the rest of the document *on this assumption*. **Keep.**
- **Scaffolding** — the author is only checking their own books (is it circular? consistent? non-duplicated?). e.g. "to avoid circularity" → the reader does not care whether the author's argument loops internally; they need the conclusion. **Drop.**

Both commonly live in the same sentence — which is why you must cut at the span level, never delete the whole sentence.

## Detection Patterns

### Observed categories (actually found in real documents)

| Category | Detection words / forms | Example (real) | Default fix mode |
|----------|-------------------------|----------------|------------------|
| A. Rhetorical stance announcement | 正直に(言うと/認め)、率直に言って、公平を期すと、誤解を恐れずに言えば、繰り返しになるが、ここで強調 / "to be honest", "frankly", "to be fair", "let me emphasize" | 「正直に認める。」 | delete / rephrase |
| B. Verification / consistency note | 循環を(避ける/回避)、整合(させる)ため、重複を(削減/排除)、矛盾しないよう / "to avoid circularity", "to keep consistent", "to avoid duplication" | 「循環の回避：…循環するため、暫定的に…」 | compress / relocate to footnote |
| D. Drafting-stage label | たたき台、ドラフト、(仮)、要推敲、TODO、書きかけ、暫定タイトル / "draft", "WIP", "tentative title" | 「仕様のたたき台：」「移行プラン（たたき台）」 | delete (turn into content if the provisional nature is substantive) |

### Speculative categories (not yet observed — enable only when needed)

| Category | Detection words / forms | Default fix mode |
|----------|-------------------------|------------------|
| C. Instruction / prompt residue | imperative instruction left verbatim in the body ("show this in a comparison table", etc.); verbatim option labels | rephrase |
| F. Traces of over-deliberation | "I went back and forth between A and B", "after considering various options" (distinguish from a legitimate lead-in that pre-empts a reader's misunderstanding) | delete (keep if legitimate) |

Per YAGNI, C and F are not treated as load-bearing until they actually appear. Enable them when observed.

## Fix Modes (section-level surgery — never delete whole sentences)

- **delete** — the span is scaffolding only.
- **compress** — scaffolding and substance share one sentence; cut only the scaffolding span and keep the substance.
- **relocate** — verification logic a reviewer needs but the reader does not; move to a footnote or appendix, not the body.
- **rephrase** — replace a stance announcement with a direct statement of fact.

## Grep Detection Patterns

Use the following with the Grep tool to scan the target file during review. High false-positive rate — pass every hit through the discrimination test above.

### Japanese patterns

```
正直に|率直に言|公平を期す|誤解を恐れ|繰り返しになるが|ここで強調|循環を避け|循環の回避|整合させ(る)?ため|重複を(削減|排除)|矛盾しないよう|たたき台|ドラフト|要推敲|書きかけ|（仮）|\(仮\)
```

### English patterns

```
to be honest|frankly|to be fair|let me emphasize|to avoid circularity|to keep .* consistent|to avoid duplication|\bdraft\b|\bWIP\b|tentative title
```

## Fix Principles (in priority order)

1. **Cut at the span level, not the sentence (most important).** Most sentences mix scaffolding with substance.
2. **Keep the substance** (the provisional value, the limitation, the conclusion). Drop only "why the author built it this way".
3. **If it bears on the reader's concern, keep it as a rephrased caveat.** If it is the author's books, drop it.
4. **If reviewer-facing rationale is genuinely needed, relocate it to a footnote or appendix** — never leave it in the body.

## Acceptable Expressions (preventing over-removal)

- **Stating that something is provisional / undecided** ("tentatively assume Read", "decided separately") = content the reader must take as a premise. Keep.
- **Pre-empting a reader's real concern** ("this looks like a cost increase, but the actual cost is…") = legitimate objection handling. Keep.
- **Process descriptions inside a changelog / appendix** (that is where process belongs). Keep.
- **Review-request comments** (`review:` etc.) are a separate channel — out of scope here.

## Delegation (no duplicate implementation)

Editing-history leaks ("formerly", "originally", "previously", "we now retract", "in the prior version" / 「かつて」「当初は」「以前は」「撤回する」「前版では」) are the *temporal* version of the same "state only current facts" principle and belong to [temporal-expressions.md](temporal-expressions.md). They are not handled here.

The two guides are complementary, not overlapping: temporal-expressions covers staleness over time (and feeds Criteria 3 and 5); process-scaffolding covers categories B and D — verification notes and drafting labels that have no temporal dimension — and feeds Criterion 3 (Writing Quality) only.

## Before/After Examples

Japanese constructions are preserved in the examples below, since the scaffolding patterns are specific to the original wording.

### Example 1: Stance announcement → delete

**Before:** 正直に認める。「影響範囲はほぼ同じ」は一般命題として成立せず、集約は…単一障害点を作る。

**After:** 「影響範囲はほぼ同じ」は一般命題として成立しない。集約は…単一障害点を作る。

### Example 2: Verification note + substance in one sentence → compress

**Before:** 循環を避けるため Base permission は暫定的に「現状の Read」…を前提に評価する。

**After:** Base permission は暫定的に「現状の Read」を前提に評価する（最終値は別途決定）。

### Example 3: Labeled verification note → relocate to footnote or delete

**Before:** 循環の回避：「背反するか」は独立要件…を通じ Base permission の値に依存する。完全に未決だと「分割は最小限」と言えず循環するため、暫定的に「現状の Read」を前提に評価する。

**After (body):** 「背反するか」の判定は Base permission の暫定値（現状の Read）を前提に行う。

**After (footnote, if needed):** ※ 最終値が未決だと判定基準が定まらないため、暫定値を置く。

### Example 4: Drafting-stage label → delete / turn into content

**Before:** ## 11. 移行プラン（たたき台）

**After:** ## 11. 移行プラン

If the provisional nature is substantively important, do not encode it as a stage label in the heading; state it as content in the body — e.g. "This section is a proposed direction; the final decision is made separately."
