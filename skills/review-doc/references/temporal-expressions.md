# Temporal Expression Detection Patterns and Fix Guide

Temporal expressions in documents become ambiguous over time and can mislead readers. This guide presents the patterns to detect and how to fix them.

## Why This Is a Problem

- **When is "recently"?** — Clear at the time of writing, but six months later the reader has no idea when it refers to
- **How long is "the latest" the latest?** — Becomes false with every version upgrade
- **"Will be addressed soon"** — Readers cannot tell whether it has been addressed or not
- **"Has been changed"** — Unclear when the change happened, so readers cannot judge whether it affects them

Temporal expressions are the biggest factor in shortening a document's "shelf life."

## Detection Patterns

### 1. Relative time expressions

| Detection pattern | Fix example |
|-------------------|-------------|
| 最近、先日、先月、来月、今後 | 2026年1月、v3.2リリース時 |
| recently, last month, soon | January 2026, in v3.2 |
| 近年、ここ数年 | 2024年以降 |
| in recent years, over the past few years | since 2024 |

### 2. Implicit recency claims

| Detection pattern | Fix example |
|-------------------|-------------|
| 最新の、現在の、現行の | v3.2（2026年1月時点） |
| the latest, the current, the new | v3.2 (as of January 2026) |
| 最新版、新しいバージョン | v3.2 |
| the newest version, the new version | v3.2 |

### 3. Vague plans and futures

| Detection pattern | Fix example |
|-------------------|-------------|
| 近いうちに、そのうち | 2026年Q2（EPT-1234） |
| 予定（日付なし）、対応予定 | 2026年3月対応予定（EPT-1234） |
| soon, in the future, upcoming | Q2 2026 (EPT-1234) |
| planned (without date), to be addressed | planned for March 2026 (EPT-1234) |

### 4. Implicit comparison expressions

| Detection pattern | Fix example |
|-------------------|-------------|
| 従来の方法、以前は | v2.x以前の方法 |
| 新しい機能、新機能 | v3.0で追加された機能 |
| previously, the old way | prior to v3.0 |
| the new feature, new functionality | feature added in v3.0 |

### 5. State-transition sentence-ending patterns

Passive or state-change sentence endings without an accompanying date or version.

**Japanese patterns:**

| Detection pattern | Fix example |
|-------------------|-------------|
| 〜された、〜されました | v3.0（2026年1月）で〜された |
| 〜になった、〜になりました | 2026年1月に〜になった |
| 〜ようになった、〜ようになりました | v3.0から〜ようになった |
| 〜ことになった、〜ことになりました | 2026年1月より〜ことになった |

**English patterns:**

| Detection pattern | Fix example |
|-------------------|-------------|
| has been changed/updated/removed | was changed in v3.0 (January 2026) |
| was deprecated | was deprecated in v3.0 (January 2026) |

These sentence-ending patterns are acceptable if a date or version is present in the same sentence.

## Grep Detection Patterns

Use the following patterns with the Grep tool to scan the target file during review.

### Japanese patterns

```
最近|先日|先月|来月|今後|近年|ここ数年|最新の|現在の|現行の|最新版|新しいバージョン|近いうちに|そのうち|対応予定|従来の|以前は|新しい機能|新機能
```

### English patterns

```
recently|last month|soon|the latest|the current|the new|in the future|upcoming|previously|the old way
```

### Sentence-ending patterns (Japanese)

```
(された|されました|になった|になりました|ようになった|ようになりました|ことになった|ことになりました)。$
```

**Note:** Sentence-ending patterns have a high false-positive rate. Cases where a date or version appears in the same sentence are acceptable. Manually verify detection results.

## Fix Principles (in priority order)

### 1. State only current facts (highest priority)

Remove temporal expressions and directly describe the current state. In most cases, the history of changes is unnecessary for the reader — the current facts are sufficient.

**Before:** The UI was recently changed. Configure settings from the new screen.

**After:** The settings screen is at Settings > Security.

### 2. Pin to a date or version

If the history of changes is important, add a specific point in time.

**Before:** The authentication method was recently changed.

**After:** The authentication method was changed to OAuth 2.0 in v3.0 (January 2026).

### 3. Delete if unnecessary

Remove statements like "will be addressed in the future" if there is no concrete plan.

**Before:** This feature will be improved soon.

**After:** (Delete, or include only if a specific plan exists)

### 4. Move to operational metadata

Place scheduled dates in an operational information section rather than the body text.

**Before:** "v4.0 is expected to be released next month" in the body text

**After:** In the operational information section: "Next update: March 2026 (after v4.0 release)"

## Acceptable Temporal Expressions

The following cases allow temporal expressions:

- **Within operational metadata**: Date-annotated expressions within the operational information section, such as "Last updated: January 2026" or "Next review: July 2026"
- **When a date is co-located**: When a specific date or version appears in the same sentence, such as "feature added in January 2026"
- **Within changelogs**: Descriptions within CHANGELOG or update history sections
- **With textlint disable comments**: Cases intentionally allowed with `<!-- textlint-disable -->`

## Before/After Examples

### Example 1: Feature description

**Before:**
> The recently added new dashboard feature lets you visualize team activity. It is significantly improved compared to the previous reporting feature.

**After:**
> The dashboard (added in v3.0) visualizes team activity. Compared to the reporting feature in v2.x and earlier, it supports real-time updates and automated aggregation.

### Example 2: Configuration procedure

**Before:**
> Due to a recent security policy change, you need to migrate to the new authentication method. All users will be affected soon.

**After:**
> Due to the security policy change (January 2026, SEC-456), migration to OAuth 2.0 authentication is required. The migration deadline for all users is March 31, 2026.

### Example 3: Troubleshooting

**Before:**
> This issue was fixed in the latest version. If you're still using an older version, please update.

**After:**
> This issue was fixed in v3.1.2. If using v3.1.1 or earlier, update to v3.1.2 or later.
