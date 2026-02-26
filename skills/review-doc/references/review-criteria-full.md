# Review Criteria (Full Version)

Detailed review criteria with checklists and grading rubrics. For a quick review, see [review-checklist.md](review-checklist.md).

## Criterion 1: Structure

Evaluate the information architecture and organization of the document.

**Checklist:**

- [ ] Is it structured as an inverted pyramid (most important information first)?
- [ ] Is the hierarchy shallow (maximum 2 levels deep)?
- [ ] Is there a table of contents or navigation?
- [ ] Do headings accurately represent the content (can you understand the overview by reading headings alone)?
- [ ] Is an appropriate content type chosen (conceptual explanation / how-to guide / reference / quickstart / troubleshooting)?
- [ ] Are there links to related documents?

**Grading Criteria:**

| Grade | Criteria |
|-------|----------|
| A | All checklist items are met; readers can find information without confusion |
| B | Most items are met, but some room for improvement |
| C | Structural issues exist that affect information findability |
| D | Structure is inadequate; readers cannot find the information they need |

## Criterion 2: Accuracy & Completeness

Evaluate technical accuracy and topic coverage.

**Checklist:**

- [ ] Is the information technically accurate?
- [ ] Is it consistent with official documentation and authoritative sources?
- [ ] Does it comprehensively cover the information needed for the topic (cross-reference with relevant guidelines if available)?
- [ ] Do code examples and configuration samples work correctly?
- [ ] Are sources and reference links provided?
- [ ] Are version information and last-updated dates stated?

**Grading Criteria:**

| Grade | Criteria |
|-------|----------|
| A | Accurate and comprehensive, with sources clearly cited |
| B | Generally accurate, but some information is missing or sources are lacking |
| C | Contains inaccurate information or important topics are missing |
| D | Contains critical errors or content is significantly incomplete |

**Note:** Grades for this criterion may be adjusted in Step 5 (WebSearch Validation)
if external sources contradict the review findings. See [web-validation.md](web-validation.md).

## Criterion 3: Writing Quality

Evaluate readability, consistency, and scannability of the writing.

**Checklist:**

- [ ] Does it use plain language (are technical terms explained)?
- [ ] Is active voice used?
- [ ] Is the one-idea-per-sentence principle followed?
- [ ] Is the writing style consistent (no mixing of formal/informal registers)?
- [ ] Is scannability high (appropriate use of bullet lists, tables, code blocks)?
- [ ] Is emphasis (bold, links) kept to 10% or less of the content?
- [ ] Are temporal expressions made specific ("recently", "currently" → specific dates/versions; refer to Step 3 detection results)?

**Grading Criteria:**

| Grade | Criteria |
|-------|----------|
| A | Readable, scannable, and consistent in style |
| B | Generally readable, but partially needs improvement |
| C | Has hard-to-read sections with style inconsistencies or verbose expressions |
| D | Overall difficult to read; requires significant rewriting |

**Note:** textlint findings (terminology consistency, verbose expressions, punctuation) primarily affect this criterion. See [textlint-integration.md](textlint-integration.md).

## Criterion 4: Practicality

Evaluate whether readers can take action based on the document.

**Checklist:**

- [ ] Is the "why" (background/rationale) explained (not just rules, but reasoning)?
- [ ] Are there concrete examples or sample code/configurations?
- [ ] Are steps clear and executable as-is by the reader?
- [ ] Are prerequisites explicitly stated?
- [ ] Are exceptions and edge cases addressed?
- [ ] Can someone who joined today understand it (no assumed specialized knowledge)?

**Grading Criteria:**

| Grade | Criteria |
|-------|----------|
| A | Specific enough to act on immediately, with sufficient background explanation |
| B | Generally practical, but some steps or examples are missing |
| C | Lacks practicality; readers need to search for additional information |
| D | Too abstract to act on |

**Note:** Grades for this criterion may be adjusted in Step 5 (WebSearch Validation)
if external sources contradict the review findings. See [web-validation.md](web-validation.md).

## Criterion 5: Maintainability

Evaluate whether the document can be maintained over the long term.

**Checklist:**

- [ ] Are external links not excessive (managing link-rot risk)?
- [ ] Is content not duplicated (linking to other documents rather than copying)?
- [ ] Is company-specific information separated from general information (general info links to official docs)?
- [ ] Is the document's content type (how-to guide vs. conceptual explanation) clearly defined, and does it delegate detailed explanations to official documentation rather than duplicating them?
- [ ] Are there no expressions that become stale over time ("the latest ..." / "new version" → specific version numbers/dates; see [temporal-expressions.md](temporal-expressions.md) for details)?

**Grading Criteria:**

| Grade | Criteria |
|-------|----------|
| A | Easy to maintain long-term, with update mechanisms in place |
| B | Generally maintainable, but some room for improvement |
| C | Maintenance concerns exist with high risk of becoming outdated |
| D | Difficult to maintain; likely to become outdated quickly |

**Note:** Grades for this criterion may be adjusted in Step 5 (WebSearch Validation)
if external sources contradict the review findings. See [web-validation.md](web-validation.md).
