# Repository Guidelines

## Project Structure & Module Organization
This repo hosts Guitar-Specs, a web app for searching, viewing, and comparing guitars. Keep it as a monorepo with clear separation between web and API layers:

- `frontend/` — Next.js app (React components, pages, hooks).
- `backend/` — Go API (handlers, services, data access).
- `shared/` — cross-platform types or utilities intended for future React Native reuse.
- `assets/` — static media (images, icons, fixtures).
- `docs/` — architecture notes and API contracts.

Example: `frontend/src/app/guitars/[id]/page.tsx` and `backend/internal/handlers/guitars.go`.

Primary datastore: PostgreSQL (self-hosted initially).
Data model source of truth: `DATABASE.md`.

## Build, Test, and Development Commands
Document the exact commands once tooling is initialized. Use npm for the frontend by default:

- `frontend/` — `npm run dev`, `npm run build`, `npm test`.
- `backend/` — `go run ./cmd/api`, `go test ./...`.

Use containers for local development and builds; see `README.md` for the full setup and command list. Infrastructure is managed with Terraform; keep module and environment conventions in `docs/infra/`.

## Coding Style & Naming Conventions
- Frontend: TypeScript + React, 2-space indentation, `PascalCase` for components, `camelCase` for hooks and helpers.
- Backend: Go fmt defaults (`gofmt`), `CamelCase` for exported identifiers.
- File names: `kebab-case` in frontend (e.g., `guitar-card.tsx`).

Use Tailwind CSS for styling; prefer utility-first classes and small reusable components.

Frontend libraries (required):
- `next` (Next.js)
- `react`
- `tailwindcss`
- Component library: `shadcn/ui`

Backend libraries (required):
- Routing: `github.com/go-chi/chi/v5`
- DB: `github.com/jackc/pgx/v5`
- SQL codegen: `github.com/sqlc-dev/sqlc`
- Migrations: `github.com/golang-migrate/migrate/v4`
- Config: `github.com/spf13/viper`
- Logging: `log/slog` (stdlib, Go 1.21+)

## Version Baseline
- Go: `1.25.x`
- Node.js: `24.x`
- Next.js: `16.x`
- React: `19.x`
- Tailwind CSS: `4.1.x`
- shadcn/ui: `latest` (via `npx shadcn-ui@latest`)
- chi: `v5`
- pgx: `v5`
- sqlc: `v1.x`
- golang-migrate: `v4`
- viper: `v1.x`

## Version References
- Go: https://go.dev/doc/
- Node.js 24: https://nodejs.org/docs/latest-v24.x/api/
- Next.js: https://nextjs.org/docs
- React: https://react.dev/
- Tailwind CSS: https://tailwindcss.com/docs
- shadcn/ui: https://ui.shadcn.com/docs
- chi v5: https://pkg.go.dev/github.com/go-chi/chi/v5
- pgx v5: https://pkg.go.dev/github.com/jackc/pgx/v5
- sqlc: https://docs.sqlc.dev/
- golang-migrate: https://github.com/golang-migrate/migrate
- viper: https://github.com/spf13/viper

## UI/UX Design Guidelines
Design should be clear and steady with a purple tone:

- Use a restrained palette (e.g., purple accents with neutral backgrounds).
- Prioritize readability and consistent spacing across pages.
- Keep comparison tables and spec lists scannable.

## SEO, Performance, and Accessibility
The site must be optimized for SEO and Core Web Vitals while following current usability and accessibility rules, including the European Accessibility Act:

- Use semantic HTML, descriptive page titles, and structured data where applicable.
- Track and optimize LCP, INP, and CLS; keep images responsive and lazy-loaded.
- Ensure keyboard navigation, visible focus states, and sufficient color contrast.

## Accessibility Checklist
- All interactive elements are keyboard operable with visible focus states.
- Provide text alternatives for images and meaningful labels for controls.
- Meet color contrast requirements for text and UI controls.
- Use proper heading order and landmark regions for navigation.
- Avoid motion that cannot be reduced; respect `prefers-reduced-motion`.

## Observability
- Log structured JSON with `request_id`, `method`, `path`, `status`, and `duration_ms`.
- Propagate `request_id` from incoming headers (or generate if missing).
- Avoid logging PII; redact secrets and tokens.
- Capture basic API metrics (request count, latency, error rate).

## CI Expectations
- Block merges on failing tests and linting.
- Run frontend unit tests (Vitest) and backend tests (`go test ./...`).
- Enforce formatting: frontend formatter/linter and `gofmt` for Go.
- Require docs update when schema or public behavior changes.

## API Conventions
- Base path: `/api/v1`.
- Use RESTful resources: `/guitars`, `/guitars/{id}`, `/brands`.
- Pagination: `?page=1&pageSize=20` and return `total`, `page`, `pageSize`.
- Filtering: repeatable `filter` params (e.g., `filter=type:electric&filter=brand:Fender`).
- Sorting: `?sort=field:asc` or `?sort=field:desc`.
- Errors: JSON with `code`, `message`, and optional `details`.

## Data & Seed Rules
- Canonical seed/test dataset is `data/guitars.json` and must match `DATABASE.md`.
- Seed scripts (when added) should live in `backend/cmd/seed` or `backend/internal/seed`.
- Any schema change requires updating `data/guitars.json` and related loaders.
- In containers, the backend auto-applies migrations and seeds when the database is empty.

## Frontend Routing & SEO (Next.js)
- Use the App Router (`frontend/src/app`) with route segments for `/(marketing)`, `/guitars`, `/guitars/[slug]`, and `/compare`.
- Generate SEO metadata via `generateMetadata` per route; include canonical URLs and Open Graph/Twitter tags.
- Prefer static generation with ISR for list/detail pages; keep dynamic routes crawlable.
- Use clean, stable slugs; avoid query-only pages for primary navigation.

## Testing Guidelines
- Frontend: Vitest + React Testing Library.
- Backend: Go’s `testing` package for unit and handler tests.
- Name tests with clear intent (e.g., `guitar-search.test.ts`, `TestGetGuitarByID`).
- Run relevant tests as the last step of each plan execution and keep tests updated with changes.
 - For changes to Docker, env, data, or migrations, run `scripts/smoke.sh`.

## Unit Test Expectations
- Backend: add unit tests for handlers, query mappers, and seed/loaders; include success and error paths.
- Frontend: add unit tests for page rendering, empty states, and key interactions.
- Any new feature or behavior change must include or update tests.
- Avoid adding unused code paths without accompanying tests.

## Commit & Pull Request Guidelines
- Use short, imperative commit messages (e.g., "Add comparison grid").
- Commit message format: `type(scope): summary` (Conventional Commits), e.g. `feat(api): add guitars list endpoint`.
- PRs must include: summary, tests run, and screenshots for UI changes.
- Link relevant issues or tasks when available.
- Keep `README.md`, `AGENTS.md`, and `DATABASE.md` updated with any relevant changes.
- Ask for explicit approval before adding any new external library.
- Provide a brief plan and wait for approval before making any code changes.
- `gofmt -w *.go` is pre-approved when needed.
- Running backend tests (`go test ./...`) and frontend tests (`npm test`) is pre-approved.
- Running `curl` for debugging is pre-approved.

## Security & Configuration Tips
Keep secrets out of the repo. Use `.env` files for local config and list required variables in `docs/config.md`.

## Deployment & CDN
The public site is hosted at `https://www.guitar-specs.com/` and uses the free tier of Cloudflare CDN. Note any cache rules or headers in `docs/deployment.md` once the stack is live.
