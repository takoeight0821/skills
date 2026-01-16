# Product Guidelines

## Prose Style
- **Clarity and Precision:** Use clear, unambiguous language. When describing technical processes (like skill synchronization or VM lifecycle), be exact about the steps and expected outcomes.
- **Direct Address:** Use the imperative mood for instructions (e.g., "Run the script," "Configure the environment") and direct address ("You") when explaining concepts.
- **Developer-Focused:** Assume the reader is a developer. Use industry-standard terminology (Git, VM, SSH, Containerization) without excessive explanation.
- **Brevity:** Keep documentation concise. Use bullet points and tables where possible to make information scannable.

## Documentation Structure
- **README-First:** The `README.md` should serve as the entry point, providing a high-level overview, prerequisites, and quick-start instructions.
- **Detailed Guides:** Complex topics (like the sync logic or Docker networking) should have dedicated markdown files in a `docs/` or `research/` directory.
- **Consistent Formatting:** Use standard Markdown headers, code blocks (with language hints), and tables to maintain a professional and readable look.

## Coding Style & Comments
- **Self-Documenting Code:** Prioritize clean, idiomatic code (especially in Go and Shell) where the intent is clear from variable and function names.
- **Strategic Commenting:** Use comments to explain the "why" rather than the "what," especially for complex logic like SSH agent forwarding detection or the skill manifest tracking.
- **Standardized Headers:** For major scripts and Go packages, provide a brief header comment describing the file's purpose and any critical dependencies.

## Visual & CLI Identity
- **CLI Output:** Ensure CLI tools (`jig`, scripts) provide clear, actionable feedback. Use standard output for information and standard error for issues.
- **Consistent Naming:** Use consistent naming conventions for `mise` tasks, environment variables, and CLI flags to ensure a predictable user experience.
