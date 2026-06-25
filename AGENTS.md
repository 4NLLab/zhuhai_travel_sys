# Repository Guidelines

## Persistent AI Rules

Before doing any work in this repository, read and follow the root `CLAUDE.md`.
Treat `CLAUDE.md` as always-on project rules, not as a one-time conversation preference.
If these repository guidelines conflict with `CLAUDE.md`, prefer `CLAUDE.md` unless a higher-priority system, developer, or explicit user instruction says otherwise.

Additional repository rule files live in `.claude/rules/`. These `.mdc` files are rule sources, not reference notes, and must be read when their frontmatter or trigger conditions match the current task. In particular, read `.claude/rules/git-submit.mdc` before any user-authorized `git commit`, commit-message writing, PR creation, GitHub push, or related GitHub operation.

For architecture questions, backend module relationships, API surface analysis, or architecture-impacting changes, consult `backend/docs/ARCHITECTURE.md` as the project architecture reference before answering or editing code.

## Project Structure & Module Organization

This repository contains a static HTML frontend and a Go/Gin backend for a Zhuhai travel platform. Root pages such as `index.html`, `ticket.html`, `admin.html`, and `island-cruise-booking.html` are standalone frontend screens. Shared images and UI assets live in `assets/`. Backend source is in `backend/`: `main.go` starts the service, `routes/` registers API paths, `handlers/` contains HTTP controllers, `middleware/` handles auth/CORS/logging, `models/` defines GORM entities, `database/` initializes MySQL, and `config/` loads environment variables. SQL migrations and database notes live in root `database/`; backend docs live in `backend/docs/`.

## Build, Test, and Development Commands

- `cd backend && go mod download`: install backend dependencies.
- `cd backend && go run .`: run the API server, defaulting to port `8080`.
- `cd backend && go build -o server .`: build the backend binary.
- `cd backend && go test ./...`: compile and run all Go tests.
- `python3 -m http.server 8000`: serve the root static HTML pages during local frontend review.
- `mysql -u root < database/001_schema.sql`: create the local MySQL schema; apply later migrations in order.

## Coding Style & Naming Conventions

Format Go code with `gofmt` before committing. Keep packages lowercase and aligned with directory names. Use exported CamelCase for shared handlers and models, for example `AdminLogin` or `ProductSchedule`. Keep JSON fields snake_case to match existing API responses. HTML/CSS files use two-space indentation, semantic class names, and CSS custom properties in `:root`.

## Testing Guidelines

There is currently no dedicated test suite. Add Go tests as `*_test.go` files next to the package under test and prefer table-driven tests for handlers, middleware, and security helpers. Run `cd backend && go test ./...` before opening a PR. For frontend changes, manually verify the affected page and capture screenshots for visible UI changes.

## Commit & Pull Request Guidelines

Recent commits use short, imperative English subjects such as `Add island cruise booking and operations` and `Update admin module menu navigation`. Follow that pattern: start with a verb, keep the subject concise, and group related changes. Pull requests should include a summary, test results, affected pages or API routes, linked issues when available, and screenshots for UI changes.

## Security & Configuration Tips

Do not commit real secrets. Use `.env.example` or `backend/.env.example` as templates, then keep local values in ignored `.env` files. Replace `JWT_SECRET`, `PAYMENT_WEBHOOK_SECRET`, and island-cruise integration credentials outside development.
