# OpenCode Agent Instructions: Go SST v3 Monorepo

Welcome, Agent. This is a Go-based SST v3 monorepo. You must follow the structural, coding, and architectural conventions outlined below.

> ⚠️ **CRITICAL:** When you encounter file reference (e.g., @.agents/rules/general.md) , use your Read tool to load it on a need-to-know-basis. They're relevant to the SPECIFIC task at hand.

---

## 📂 Project Structure

- `internal/` - Contains all workspace packages (Go modules, web apps, etc.)
  - `internal/handlers/` - Handlers serving the web content.
- `templates/` - Html template for static content.

---

## 🛠 Code Standards & Patterns

### Go (Golang)

- Ensure all Go code is cleanly formatted (`go fmt`) and passes strict linting rules.
- **Handlers:** Keep Lambda functions inside `internal/handlers/` thin.
- **Shared Logic:** All reusable domain logic must reside in `packages/utils/`.

### Development Guidelines

For UI related code style and best practices: @./agents/rules/ui.md.

---

## 📦 Monorepo Conventions

- **Go Imports:** Import internal shared modules using the workspace's fully qualified Go module path.
  - _Example:_ `import "github.com/username/repo/packages/core/db"`
