package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"voicespin/backend/internal/conversation"
)

func newTestHandler() http.Handler {
	return NewServer(conversation.NewStore(conversation.Seed())).Routes()
}

func doRequest(t *testing.T, handler http.Handler, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	var reader *strings.Reader
	if body == "" {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func decodeConversations(t *testing.T, rec *httptest.ResponseRecorder) []conversation.Conversation {
	t.Helper()

	var conversations []conversation.Conversation
	if err := json.NewDecoder(rec.Body).Decode(&conversations); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	return conversations
}

func TestListConversations(t *testing.T) {
	handler := newTestHandler()

	tests := []struct {
		name      string
		target    string
		wantCount int
		wantIDs   []string
	}{
		{name: "all conversations", target: "/api/conversations", wantCount: 8},
		{name: "filter by status", target: "/api/conversations?status=OPEN", wantIDs: []string{"conv-001", "conv-002", "conv-004"}},
		{name: "filter by priority", target: "/api/conversations?priority=HIGH", wantIDs: []string{"conv-001", "conv-002", "conv-006"}},
		{name: "search is case-insensitive", target: "/api/conversations?search=JOHN", wantIDs: []string{"conv-001", "conv-003"}},
		{name: "search matches email", target: "/api/conversations?search=acme.io", wantIDs: []string{"conv-003"}},
		{name: "combined filters", target: "/api/conversations?status=RESOLVED&priority=LOW", wantIDs: []string{"conv-005", "conv-008"}},
		{name: "no matches returns empty list", target: "/api/conversations?search=nobody", wantCount: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, handler, http.MethodGet, tt.target, "")

			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
			}
			got := decodeConversations(t, rec)
			wantCount := tt.wantCount
			if tt.wantIDs != nil {
				wantCount = len(tt.wantIDs)
			}
			if len(got) != wantCount {
				t.Fatalf("got %d conversations, want %d: %+v", len(got), wantCount, got)
			}
			for i, id := range tt.wantIDs {
				if got[i].ID != id {
					t.Errorf("result[%d].ID = %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}

func TestListConversationsRejectsInvalidFilter(t *testing.T) {
	handler := newTestHandler()

	for _, target := range []string{
		"/api/conversations?status=ARCHIVED",
		"/api/conversations?priority=URGENT",
	} {
		rec := doRequest(t, handler, http.MethodGet, target, "")
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status = %d, want %d", target, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestGetConversation(t *testing.T) {
	handler := newTestHandler()

	rec := doRequest(t, handler, http.MethodGet, "/api/conversations/conv-001", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var c conversation.Conversation
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if c.ID != "conv-001" || c.CustomerName != "John Carter" {
		t.Errorf("unexpected conversation: %+v", c)
	}

	rec = doRequest(t, handler, http.MethodGet, "/api/conversations/missing", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("missing id: status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestPatchConversation(t *testing.T) {
	handler := newTestHandler()

	rec := doRequest(t, handler, http.MethodPatch, "/api/conversations/conv-001", `{"status":"RESOLVED","priority":"LOW"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (body: %s)", rec.Code, http.StatusOK, rec.Body.String())
	}

	var updated conversation.Conversation
	if err := json.NewDecoder(rec.Body).Decode(&updated); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if updated.Status != conversation.StatusResolved || updated.Priority != conversation.PriorityLow {
		t.Errorf("unexpected updated conversation: %+v", updated)
	}

	rec = doRequest(t, handler, http.MethodGet, "/api/conversations/conv-001", "")
	persisted := decodeOne(t, rec)
	if persisted.Status != conversation.StatusResolved || persisted.Priority != conversation.PriorityLow {
		t.Errorf("change was not persisted: %+v", persisted)
	}

	rec = doRequest(t, handler, http.MethodPatch, "/api/conversations/conv-002", `{"status":"IN_PROGRESS"}`)
	partial := decodeOne(t, rec)
	if partial.Status != conversation.StatusInProgress {
		t.Errorf("Status = %q, want %q", partial.Status, conversation.StatusInProgress)
	}
	if partial.Priority != conversation.PriorityHigh {
		t.Errorf("Priority = %q, want unchanged %q", partial.Priority, conversation.PriorityHigh)
	}
}

func TestPatchConversationValidation(t *testing.T) {
	handler := newTestHandler()

	tests := []struct {
		name       string
		target     string
		body       string
		wantStatus int
	}{
		{name: "invalid status value", target: "/api/conversations/conv-001", body: `{"status":"CLOSED"}`, wantStatus: http.StatusBadRequest},
		{name: "invalid priority value", target: "/api/conversations/conv-001", body: `{"priority":"URGENT"}`, wantStatus: http.StatusBadRequest},
		{name: "malformed JSON", target: "/api/conversations/conv-001", body: `{`, wantStatus: http.StatusBadRequest},
		{name: "no updatable fields", target: "/api/conversations/conv-001", body: `{}`, wantStatus: http.StatusBadRequest},
		{name: "unknown field rejected", target: "/api/conversations/conv-001", body: `{"id":"hacked"}`, wantStatus: http.StatusBadRequest},
		{name: "missing conversation", target: "/api/conversations/missing", body: `{"status":"OPEN"}`, wantStatus: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, handler, http.MethodPatch, tt.target, tt.body)
			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
		})
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := newTestHandler()

	rec := doRequest(t, handler, http.MethodOptions, "/api/conversations/conv-001", "")
	if rec.Code != http.StatusNoContent {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Methods"); !strings.Contains(got, "PATCH") {
		t.Errorf("Allow-Methods = %q, want it to include PATCH", got)
	}
}

func decodeOne(t *testing.T, rec *httptest.ResponseRecorder) conversation.Conversation {
	t.Helper()

	var c conversation.Conversation
	if err := json.NewDecoder(rec.Body).Decode(&c); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	return c
}
