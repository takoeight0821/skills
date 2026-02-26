# WebSearch Validation Workflow

Cross-check review findings against authoritative external sources using WebSearch and WebFetch. This workflow validates Criteria 2 (Accuracy), 4 (Practicality), and 5 (Maintainability). Criteria 1 (Structure) and 3 (Writing Quality) are not subject to external validation.

**Prerequisites:** This workflow requires access to WebSearch and WebFetch tools. If these tools are not available in the current environment, skip the entire validation step and note "External validation tools unavailable" in the output.

## Phase 1: Extract Verifiable Claims

Extract up to **10 verifiable claims** from the document. Target the following types:

- Technology names and version numbers
- API usage, configuration values, and parameters
- External URLs (links to documentation, tools, services)
- Best practice assertions
- Command examples and code snippets

**Selection priority** (highest first):

1. Items the reviewer flagged as uncertain or had low confidence about during Step 4
2. Version numbers and technology references
3. Configuration values and code examples
4. External URLs
5. Best practice assertions

## Phase 2: Search Query Construction and Execution

Construct targeted queries for each claim. Use the following templates:

| Claim Type | Tool | Query Template |
|------------|------|---------------|
| Version confirmation | WebSearch | `"<technology> latest version LTS"` |
| API / configuration | WebSearch | `"<technology> <param> official documentation"` |
| Best practice | WebSearch | `"<technology> best practices <topic>"` |
| Link verification | WebFetch | `<URL>` |

**Limits:**

- Maximum **5 WebSearch** calls
- Maximum **3 WebFetch** calls

**Skip conditions — skip this entire step and note the reason if:**

- The document is internal-only with no publicly referenceable claims
- No verifiable claims were extracted in Phase 1

## Phase 3: Cross-Check

Classify each verified claim into one of the following categories:

| Classification | Description |
|----------------|-------------|
| **Confirmed** | External source agrees with the document |
| **Outdated** | Information was once correct but a newer version/approach exists |
| **Contradicted** | External source directly contradicts the document |
| **Broken Link** | URL returns 404 or is unreachable |
| **Stale Link** | URL is reachable but content has changed significantly or is no longer relevant |
| **Unverified** | Could not find authoritative external source to confirm or deny |

## Phase 4: Scoring and Grade Adjustment

### Validation Confidence

Assign a **Validation Confidence** level to each applicable criterion (2, 4, 5):

| Confidence | Condition |
|------------|-----------|
| **High** | All checked claims confirmed; no contradictions |
| **Medium** | 1-2 minor issues (Outdated or Unverified) found |
| **Low** | Contradictions found, or multiple Outdated items |

### Grade Adjustment Rules

| Finding | Action |
|---------|--------|
| **Contradicted** found | Downgrade the affected criterion by 1 level (e.g., B -> C) |
| **Outdated** found | Add a warning in feedback; do NOT auto-downgrade |
| **Confirmed** found | Does NOT upgrade an existing low grade |
| **Broken Link** / **Stale Link** found | Add a warning; downgrade Maintainability (Criterion 5) by 1 level if 2+ broken/stale links |

## Output Format

Include the following in the review output:

### Summary Table Addition

Add a `Validation` column to the summary table:

```
| Criterion | Grade | Validation | Comment |
```

- Use `--` for Criteria 1 and 3 (not validated)
- Use `High` / `Medium` / `Low` for Criteria 2, 4, and 5

### External Validation Section

Add after the Detailed Feedback section:

```markdown
### External Validation

| # | Claim | Type | Classification | Source |
|---|-------|------|----------------|--------|
| 1 | ... | Version | Confirmed | [link] |
| 2 | ... | URL | Broken Link | [link] |

**Grade adjustments:**
- (List any grade changes and reasons, or "None")

**Skipped / Notes:**
- (If validation was skipped, state the reason here)
```
