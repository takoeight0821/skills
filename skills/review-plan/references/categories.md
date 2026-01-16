# Review Comment Categories

Detailed taxonomy for classifying code review comments.

## Bug Category

Issues that cause incorrect behavior or crashes.

**Keywords**: crash, error, null, exception, wrong, incorrect, broken, fails, undefined, NaN

**Examples**:
```javascript
// Review: This will throw if user is null
const name = user.name;

// Review: Off-by-one error in loop bounds
for (let i = 0; i <= arr.length; i++) { ... }

// Review: Race condition - data might not be loaded yet
console.log(this.data.value);
```

**Priority**: High - Fix before merging

## Security Category

Vulnerabilities and security risks.

**Keywords**: injection, XSS, CSRF, auth, permission, vulnerable, sanitize, escape, token, secret

**Examples**:
```javascript
// Review: SQL injection vulnerability
const query = `SELECT * FROM users WHERE id = ${userId}`;

// Review: XSS - user input not sanitized
element.innerHTML = userComment;

// Review: Sensitive data exposed in logs
console.log('User token:', authToken);
```

**Priority**: High - Must fix before deployment

## Performance Category

Performance issues and optimization opportunities.

**Keywords**: slow, N+1, loop, memory, cache, optimize, batch, lazy, eager, index

**Examples**:
```python
# Review: N+1 query - fetching comments in loop
for post in posts:
    comments = db.query(Comment).filter(post_id=post.id).all()

# Review: Loading entire file into memory
data = open('large_file.csv').read()

# Review: Unnecessary re-renders on every keystroke
onChange={(e) => setQuery(e.target.value)}
```

**Priority**: Medium - Address before release

## Design Category

Architectural and design improvements.

**Keywords**: refactor, extract, separate, coupling, responsibility, abstraction, duplicate, DRY, SOLID

**Examples**:
```typescript
// Review: This function does too many things - extract validation logic
function createUser(data) {
    // validation
    // database insert
    // send email
    // update cache
}

// Review: Tight coupling to external service - introduce interface
const result = await StripeAPI.charge(amount);

// Review: Duplicate logic - see similar code in OrderService
```

**Priority**: Medium - Improve maintainability

## Question Category

Clarifications and understanding.

**Keywords**: why, what, how, unclear, explain, purpose, reason, intention, confused

**Examples**:
```go
// Review: Why is this hardcoded to 42?
const maxRetries = 42

// Review: What happens if this returns empty?
results := fetchData()

// Review: Unclear variable name - what does 'x' represent?
x := calculate(a, b)
```

**Priority**: Low - Clarify before final review

## Style Category

Code style and conventions.

**Keywords**: naming, format, convention, consistency, typo, comment, documentation

**Examples**:
```java
// Review: Method name should be camelCase
public void Process_Data() { ... }

// Review: Missing JSDoc for public API
export function transformData(input) { ... }

// Review: Inconsistent indentation
if (condition) {
  doA();
    doB();
}
```

**Priority**: Low - Nice to fix but not blocking

## Classification Algorithm

When categorizing a review comment:

1. **Check for security keywords first** - Security issues are always high priority
2. **Check for bug keywords** - Bugs need immediate attention
3. **Check for performance keywords** - Performance issues affect user experience
4. **Check for design keywords** - Design improvements help long-term maintenance
5. **Check for question marks or question words** - Questions need answers before proceeding
6. **Default to style** - Remaining issues are likely style-related

## Priority Matrix

| Category | Default Priority | Can Escalate To |
|----------|------------------|-----------------|
| Security | High | Critical |
| Bug | High | Critical |
| Performance | Medium | High |
| Design | Medium | High |
| Question | Low | Medium |
| Style | Low | Low |

## Reviewer Attribution

When comments include reviewer name `Review(username):`:
- Track who raised each concern
- Group related comments by reviewer
- Consider reaching out for clarification on ambiguous feedback
