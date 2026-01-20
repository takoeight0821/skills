---
name: review-plan
description: This skill should be used when the user invokes "/review-plan", asks to "レビューを読んで", "レビュー対応", "レビュー指摘を修正", "改善計画を立てて", or wants to address code review comments. Searches for "review:" or "review(username):" comments in code and creates an improvement plan in plan mode.
---

# Review Plan

Read human code review comments embedded in source code and create a structured improvement plan.

## Overview

This skill processes code review feedback written directly in source files as comments. review comments follow the pattern `review:` or `review(username):` and contain feedback about bugs, design improvements, questions, or performance concerns.

## Comment Format

Review comments are embedded in source code using standard comment syntax:

```
// review: This could throw NullPointerException
# review(john): Consider extracting this to a separate function
/* review: N+1 query issue here */
-- review: Type signature could be more specific
```

The format supports:

- `review:` - Anonymous review comment
- `review(username):` - review comment with reviewer attribution

## Workflow

### 1. Search for review Comments

Use Grep to find all review comments in the codebase:

```
Grep: pattern="review(\([^)]*\))?:"
```

This matches both `review:` and `review(username):` patterns.

### 2. Collect and Read Context

For each found comment, read the surrounding code to understand the context:

- What file and function contains the comment
- What code the review is referring to
- Any related code that might be affected by changes

### 3. Categorize Comments

Classify each review comment into categories:

| Category    | Keywords                                              | Priority |
| ----------- | ----------------------------------------------------- | -------- |
| Bug         | crash, error, null, exception, wrong, incorrect       | High     |
| Security    | injection, XSS, auth, permission, vulnerable          | High     |
| Performance | slow, N+1, loop, memory, cache                        | Medium   |
| Design      | refactor, extract, separate, coupling, responsibility | Medium   |
| Question    | why, what, how, unclear, explain                      | Low      |
| Style       | naming, format, convention, consistency               | Low      |

### 4. Prioritize Issues

Order issues by:

1. **Severity**: Bugs and security issues first
2. **Impact**: Changes affecting multiple files/functions
3. **Dependencies**: Issues that block other fixes

### 5. Enter Plan Mode

Call EnterPlanMode to create a structured improvement plan. The plan should include:

1. **Summary of review findings**
   - Total number of comments
   - Breakdown by category
   - Files affected

2. **Prioritized action items**
   - Each issue with clear description
   - Proposed solution approach
   - Files to modify

3. **Implementation order**
   - Group related changes
   - Consider dependencies between fixes

## Output Structure

When in plan mode, structure the plan as:

```markdown
# review Response Plan

## Summary

- Total comments: N
- Categories: X bugs, Y design, Z questions...
- Files affected: list

## High Priority

### 1. [Bug] Issue title

- **File**: path/to/file.ts:42
- **Comment**: Original review comment
- **Proposed fix**: Description of solution
- **Impact**: What else might be affected

### 2. [Security] Issue title

...

## Medium Priority

### 3. [Design] Issue title

...

## Low Priority

### 4. [Question] Clarification needed

...

## Implementation Order

1. Fix X first (blocks Y and Z)
2. Then address Y
3. Finally Z
```

## Tips

### Writing Effective review Comments

When leaving review comments for this skill to process:

- Be specific about the location and issue
- Suggest concrete improvements when possible
- Use category keywords for automatic classification:
  - "Bug:", "Error:", "Crash:" for bugs
  - "Security:", "Vulnerable:" for security
  - "Slow:", "Performance:" for performance
  - "Refactor:", "Extract:" for design
  - "Why?", "Unclear:" for questions

### Handling Large Codebases

For large codebases with many review comments:

- Search specific directories: `Grep pattern in src/`
- Filter by file type: `Grep pattern glob="*.ts"`
- Process incrementally by priority

### After Planning

Once the plan is approved:

1. Work through items in priority order
2. Remove or update review comments as issues are addressed
3. Run tests after each significant change

## Additional Resources

For detailed category definitions and examples:

- **`references/categories.md`** - Complete category taxonomy with examples
