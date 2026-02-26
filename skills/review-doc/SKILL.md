---
name: review-doc
description: >
  This skill should be used when the user asks to review, check, proofread, improve,
  or give feedback on a document or Markdown file. Evaluates structure, accuracy,
  writing quality, practicality, and maintainability using a five-criteria grading system.
  Triggers: "review doc", "review this document", "review this markdown",
  "check my document", "check this README", "proofread", "give feedback on this doc",
  "document review", "improve my documentation", "docs review",
  "ドキュメントをレビュー", "レビューして", "ドキュメントのチェック", "文章レビュー".
---

# Review Doc — Document Review

A skill that reviews Markdown documents across five criteria: Structure, Accuracy, Writing Quality, Practicality, and Maintainability.

## Workflow

### Step 1: Identify the review target

1. If the user specified a file path, use that file
2. If not specified, list Markdown files under the current directory and ask which one to review
3. Read the target file

### Step 2: Gather applicable guidelines

Identify relevant style guides and project conventions for the review:

- Refer to any relevant guidelines for the target document (e.g., README.md, CONTRIBUTING.md, style guides in the same directory or project root)
- Load the review criteria from [references/review-criteria-full.md](references/review-criteria-full.md)
- If the project has textlint configured, follow [references/textlint-integration.md](references/textlint-integration.md) to run mechanical checks and reflect findings in the corresponding criteria

### Step 3: Scan for temporal expressions

Use the Grep tool to scan the target file with the regex patterns listed in [references/temporal-expressions.md](references/temporal-expressions.md). Reflect the detection results in the Writing Quality (Criterion 3) and Maintainability (Criterion 5) evaluations in Step 4.

### Step 4: Review across five criteria

Evaluate the target document across the five criteria below. For a quick review, use [references/review-checklist.md](references/review-checklist.md) instead of the full criteria. For each criterion, assign a grade of **A (Excellent) / B (Good) / C (Needs Improvement) / D (Requires Revision)** and provide specific feedback. Reflect the Step 3 detection results in Criteria 3 and 5.

### Step 5: Output the review results

Output the review results in the following format:

```markdown
## Review Results: [filename]

### Summary

| Criterion | Grade | Comment |
|-----------|-------|---------|
| Structure | [grade] | ... |
| Accuracy & Completeness | [grade] | ... |
| Writing Quality | [grade] | ... |
| Practicality | [grade] | ... |
| Maintainability | [grade] | ... |

### Detailed Feedback

#### 1. Structure
...(specific findings and improvement suggestions)

#### 2. Accuracy & Completeness
...

#### 3. Writing Quality
...

#### 4. Practicality
...

#### 5. Maintainability
...

### Prioritized Improvement Actions

1. **[High]** ...
2. **[Medium]** ...
3. **[Low]** ...
```

---

## Review Criteria (Summary)

| Criterion | Focus | Key Question |
|-----------|-------|-------------|
| 1. Structure | Information architecture, navigation | Can readers find what they need? |
| 2. Accuracy & Completeness | Technical correctness, coverage | Is the information correct and complete? |
| 3. Writing Quality | Readability, consistency, scannability | Is it easy to read and scan? |
| 4. Practicality | Actionability, examples, prerequisites | Can readers act on this immediately? |
| 5. Maintainability | Long-term upkeep, link-rot, staleness | Will this stay accurate over time? |

Each criterion is graded **A (Excellent) / B (Good) / C (Needs Improvement) / D (Requires Revision)**.

For full checklists and grading rubrics, see [references/review-criteria-full.md](references/review-criteria-full.md).

---

## Review Mindset

- **Be constructive**: Always suggest specific improvements, not just point out problems
- **Prioritize**: Not everything needs to be fixed at once. Assign High / Medium / Low priorities
- **Acknowledge good points**: Mention what's well done, not just areas for improvement
- **Think from the reader's perspective**: Judge by value to the intended reader, not the reviewer's personal preferences

## References

- [references/review-criteria-full.md](references/review-criteria-full.md) — Full review criteria with checklists and grading rubrics
- [references/review-checklist.md](references/review-checklist.md) — Quick review checklist (for quick reviews)
- [references/temporal-expressions.md](references/temporal-expressions.md) — Temporal expression detection patterns and fix guide
- [references/textlint-integration.md](references/textlint-integration.md) — textlint configuration detection, execution, and setup suggestions
