# Conversation Dashboard — VoiceSpin Full Stack Assignment

A small customer support dashboard: a Go REST API (in-memory store, seed data) and an Angular 22 frontend that lists conversations, filters them, and lets you update status and priority.

## Project structure

```
backend/
  main.go                  # server wiring, listens on :8080
  internal/api/            # HTTP handlers + handler tests
  internal/conversation/   # model, in-memory store, seed data + tests
  noai/                    # standalone no-AI exercise + tests
frontend/
  src/app/core/            # conversation model + HTTP service
  src/app/features/        # conversation list and detail components
  proxy.conf.json          # dev proxy: /api → http://localhost:8080
```

## Prerequisites

- Go 1.27+ (as declared in `backend/go.mod`; uses `net/http` method-pattern routing)
- Node.js 22.22.3+ / 24.15+ / 26+ (Angular 22 requirement), npm

## Run the backend

```bash
cd backend
go run .
```

Serves on `http://localhost:8080` (override with the `PORT` environment variable). Data lives in memory and resets on restart.

## Run the frontend

```bash
cd frontend
npm install
npm start
```

Opens on `http://localhost:4200`. The dev server proxies `/api` to the backend on port 8080 (see `proxy.conf.json`), so run both at the same time.

## Run the tests

Backend (includes the no-AI exercise tests in `backend/noai`):

```bash
cd backend
go test ./...

# also useful:
go vet ./...
go test -race ./...
```

Frontend:

```bash
cd frontend
npm test -- --watch=false
npm run build
```

## API

| Method | Path | Notes |
| ------ | ---- | ----- |
| GET | `/api/conversations` | Supports `?status=OPEN`, `?priority=HIGH`, `?search=john` (case-insensitive, matches name, email, subject) |
| GET | `/api/conversations/:id` | 404 if not found |
| PATCH | `/api/conversations/:id` | Body: `{"status": "...", "priority": "..."}` — either or both fields |

## Design decisions

- **Go standard library only** — `net/http` method patterns instead of a routing framework; no external dependencies.
- **In-memory store behind a `sync.RWMutex`** — safe for concurrent requests; seed data is copied so callers cannot mutate it.
- **Strict PATCH** — only `status` and `priority` are accepted; unknown fields are rejected (`id` and `createdAt` are immutable) and an empty patch returns 400.
- **Server-side validation** — invalid filter values return 400 instead of silently returning nothing.
- **Angular 22 standalone components + signals** with `OnPush`, a thin HTTP service layer, and explicit loading / error / empty states.
- **Tests chosen for behavior, not coverage** — filtering/search rules, PATCH validation and persistence, update semantics, and the main UI flows.

## Time spent

- Total time spent: **~2.5–3 hours**
- Approximate time spent using AI: **~45 minutes**
- Approximate time spent on the no-AI section: **~45 minutes**

## What I would improve with more time

- **Debounce the search input** — currently every keystroke fires a request; fine for 8 seed rows, wrong at scale.
- **Optimistic updates** for status/priority with rollback if the PATCH fails.
- **Persistence** — swap the in-memory store for SQLite with migrations, and add the create/delete endpoints the spec left out.
- **Pagination and sorting** on the list endpoint for larger datasets.
- **End-to-end tests** (Playwright) through the dev proxy, plus CI that runs both test suites.
