
# OpenCode Agent Instructions: Go SST v3 Monorepo

Welcome, Agent. This is a Go-based SST v3 monorepo. You must follow the structural, coding, and architectural conventions outlined below.

> ⚠️ **CRITICAL:** When you encounter file reference (e.g., @.agents/rules/general.md) , use your Read tool to load it on a need-to-know-basis. They're relevant to the SPECIFIC task at hand.

---

## 📂 Project Structure

- `internal/` - Contains all workspace packages (Go modules, web apps, etc.)
  - `internal/handlers/` - Handlers serving the web content.
- `templates/` - Html template for static content.
- `pkg/utils/` - Utils commoly used in the project 
- `config/` - Contains all go environment variables needed to be loaded with os.getenv etc. Also if we have appConfig it is here. 

---

## 🛠 Code Standards & Patterns

### Go (Golang)

- Ensure all Go code is cleanly formatted (`go fmt`) and passes strict linting rules.
- **Handlers:** Keep Lambda functions inside `internal/handlers/` thin.
- **Shared Logic:** All reusable domain logic must reside in `pkg/utils/`.
- **SDK Used:** This project must use Go MCP SDK https://github.com/modelcontextprotocol/go-sdk

---

## 📦 Monorepo Conventions

- **Go Imports:** Import internal shared modules using the workspace's fully qualified Go module path.
  - _Example:_ `import "github.com/username/repo/packages/core/db"`
