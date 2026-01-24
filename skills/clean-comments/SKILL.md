---
name: clean-comments
description: Optimize code comments by removing redundant ones and keeping only valuable comments. Use when asked to "clean up comments", "remove redundant comments", "optimize comments", or after writing code that needs comment review.
---

# Clean Comments

Remove redundant comments and keep only valuable ones.

## Comments to Remove

1. **Restating code**
   ```go
   // increment i
   i++

   // get user
   user := getUser(id)
   ```

2. **Obvious type/variable explanations**
   ```go
   // user ID
   userID := 123

   // check error
   if err != nil {
   ```

3. **Empty TODO/FIXME promises**
   ```go
   // TODO: optimize later
   // FIXME: fix someday
   ```

4. **Change history** (use Git instead)
   ```go
   // 2024-01-01 John: bug fix
   // added in v2.0
   ```

5. **Commented-out code**
   ```go
   // oldFunction()
   // if debug { log() }
   ```

6. **Over-detailed implementation explanations**
   ```go
   // This loop iterates through each element in the array,
   // checks if it matches the condition, and if so,
   // appends it to the result slice.
   for _, item := range items {
       if condition(item) {
           result = append(result, item)
       }
   }
   ```

## Comments to Keep

1. **WHY explanations**
   ```go
   // Add 100ms delay to avoid rate limiting
   time.Sleep(100 * time.Millisecond)
   ```

2. **Non-obvious business logic**
   ```go
   // Japanese holidays have substitute holidays when falling on Sunday
   if holiday.Weekday() == time.Sunday {
       holidays = append(holidays, holiday.AddDate(0, 0, 1))
   }
   ```

3. **External constraints/dependencies**
   ```go
   // AWS S3 multipart upload requires minimum 5MB parts
   const minPartSize = 5 * 1024 * 1024
   ```

4. **Warnings**
   ```go
   // WARNING: changing this order causes deadlock
   mu1.Lock()
   mu2.Lock()
   ```

5. **Performance reasons**
   ```go
   // Batch fetch to avoid N+1 queries
   users := db.PreloadAll("posts").Find(&users)
   ```

6. **Public API documentation**
   ```go
   // Client wraps HTTP client with automatic retry and rate limiting.
   type Client struct {
   ```

## Decision Criteria

| Question | Yes → Keep | No → Remove |
|----------|------------|-------------|
| Explains WHY not obvious from code? | ✓ | |
| Documents external system constraints? | ✓ | |
| Important warning for future developers? | ✓ | |
| Public API documentation? | ✓ | |
| Restates what code does? | | ✓ |
| Obvious from variable/function names? | | ✓ |

## Execution

1. Read target file
2. Evaluate each comment against criteria
3. Identify comments to remove
4. Remove with Edit tool
5. Report summary of removed comments
