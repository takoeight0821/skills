# textlint Integration

General writing quality rules (terminology consistency, verbose expressions, punctuation, etc.) can be mechanically checked with textlint. Since temporal expression detection is handled by the Grep tool, textlint's role is limited to other writing rules.

## How textlint findings map to review criteria

textlint violations primarily affect **Criterion 3: Writing Quality** (terminology consistency, verbose expressions, punctuation, sentence length). Some rules may also relate to **Criterion 5: Maintainability** (e.g., link validation rules). Reflect textlint results in the corresponding criterion grades.

## Investigate the project's textlint configuration

1. Use the Glob/Grep tools to investigate textlint configuration in the project (`.textlintrc*`, `package.json` scripts, `mise.toml` task definitions, etc.)
2. If a configuration is found, run textlint following that project's method

## If textlint is not set up

If no textlint configuration is found, suggest introducing it at the end of the review results. Example:

```bash
# Initial setup
pnpm dlx textlint --init

# Preset for Japanese technical documents
pnpm add -D textlint textlint-rule-preset-ja-technical-writing

# One-off execution (no installation needed)
pnpm dlx textlint --rule preset-ja-technical-writing '**/*.md'
```
