# AI Usage

**Tool:** opencode (a terminal coding agent).
**My loop:** plan first → the agent writes code → I read every change → run the tests and the app. The agent does the typing; the plan, the API behavior, and what gets accepted are my decisions.

In my normal opencode setup I use reusable skills (planning, code review, verification). On this assignment I kept the process direct — no skills — so every change was easy to check.

## Example 1 — Planning the backend, then making the AI follow the plan

- **Goal:** Build the Go API the way I planned it, not however the AI preferred.
- **Prompt:** "Start with the Go API + Angular app + tests." With constraints: Go standard library only, in-memory data, the three endpoints from the spec, filters for status, priority, and search.
- **Result:** The AI suggested a folder split — `internal/conversation` (data) and `internal/api` (HTTP) — and wrote the handlers with a mutex-protected store, plus tests.
- **Decision:** Kept the structure. Said no to adding a routing library. Checked the PATCH behavior and kept it strict: it only accepts `status` and `priority`; unknown fields are rejected with 400, so `id` and `createdAt` can never change. Verified with `go test -race ./...`.

## Example 2 — The AI was wrong; the tests caught it before it shipped

- **Goal:** The detail panel dropdowns must show the conversation's real status and priority.
- **Prompt:** The Angular part of the same instruction: an app that consumes the API, with tests.
- **Result:** The AI wrote `<select [value]="priority()">`. The test failed: the dropdown showed the first option (LOW) instead of HIGH. Reason: Angular applies `[value]` before `@for` has created the options, so the binding silently does nothing.
- **Decision:** Threw away the version that only looked correct, kept the real fix (bind `[selected]` on each option), and kept the failing test so the bug can never come back. Lesson: AI UI code can look right and still be wrong — tests catch it.

## Example 3 — The test was wrong, not the code

- **Goal:** Test the combined status + priority filter on the list endpoint.
- **Prompt:** Same instruction — tests for the backend.
- **Result:** A test expected `RESOLVED` + `LOW` to return 1 conversation; the seed data has 2.
- **Decision:** Checked the data, confirmed the handler was right, and fixed the test instead of changing working code to match a wrong expectation.

**The rule across all three examples: nothing is accepted until it has been read, run, and proven by a test.**
