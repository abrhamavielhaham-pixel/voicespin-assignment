package conversation

import (
	"errors"
	"testing"
	"time"
)

func testSeed() []Conversation {
	return []Conversation{
		{
			ID:            "1",
			CustomerName:  "John Carter",
			CustomerEmail: "john.carter@example.com",
			Subject:       "Cannot reset my password",
			Status:        StatusOpen,
			Priority:      PriorityHigh,
			CreatedAt:     time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:            "2",
			CustomerName:  "Maria Lopez",
			CustomerEmail: "maria.lopez@example.com",
			Subject:       "Invoice charged twice",
			Status:        StatusOpen,
			Priority:      PriorityLow,
			CreatedAt:     time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:            "3",
			CustomerName:  "Emily Watson",
			CustomerEmail: "emily@acme.io",
			Subject:       "SSO login fails",
			Status:        StatusResolved,
			Priority:      PriorityHigh,
			CreatedAt:     time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC),
		},
	}
}

func TestListFilters(t *testing.T) {
	store := NewStore(testSeed())

	tests := []struct {
		name   string
		filter Filter
		want   []string
	}{
		{name: "no filter returns all", filter: Filter{}, want: []string{"1", "2", "3"}},
		{name: "filter by status", filter: Filter{Status: StatusOpen}, want: []string{"1", "2"}},
		{name: "filter by priority", filter: Filter{Priority: PriorityHigh}, want: []string{"1", "3"}},
		{name: "search matches name case-insensitively", filter: Filter{Search: "JOHN"}, want: []string{"1"}},
		{name: "search matches email", filter: Filter{Search: "acme.io"}, want: []string{"3"}},
		{name: "search matches subject", filter: Filter{Search: "invoice"}, want: []string{"2"}},
		{name: "combined status and priority", filter: Filter{Status: StatusOpen, Priority: PriorityHigh}, want: []string{"1"}},
		{name: "combined filters with search", filter: Filter{Status: StatusResolved, Search: "emily"}, want: []string{"3"}},
		{name: "no matches returns empty slice", filter: Filter{Search: "nobody"}, want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store.List(tt.filter)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d conversations, want %d: %+v", len(got), len(tt.want), got)
			}
			for i, id := range tt.want {
				if got[i].ID != id {
					t.Errorf("result[%d].ID = %q, want %q", i, got[i].ID, id)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	store := NewStore(testSeed())

	c, err := store.Get("2")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if c.CustomerName != "Maria Lopez" {
		t.Errorf("CustomerName = %q, want %q", c.CustomerName, "Maria Lopez")
	}

	if _, err := store.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(missing) error = %v, want ErrNotFound", err)
	}
}

func TestUpdate(t *testing.T) {
	store := NewStore(testSeed())
	status := StatusResolved
	priority := PriorityMedium

	updated, err := store.Update("1", Update{Status: &status})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Status != StatusResolved {
		t.Errorf("Status = %q, want %q", updated.Status, StatusResolved)
	}
	if updated.Priority != PriorityHigh {
		t.Errorf("Priority = %q, want unchanged %q", updated.Priority, PriorityHigh)
	}
	if updated.CustomerName != "John Carter" {
		t.Errorf("CustomerName = %q, want unchanged %q", updated.CustomerName, "John Carter")
	}

	updated, err = store.Update("1", Update{Priority: &priority})
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if updated.Priority != PriorityMedium {
		t.Errorf("Priority = %q, want %q", updated.Priority, PriorityMedium)
	}
	if updated.Status != StatusResolved {
		t.Errorf("Status = %q, want unchanged %q", updated.Status, StatusResolved)
	}

	if _, err := store.Update("missing", Update{Status: &status}); !errors.Is(err, ErrNotFound) {
		t.Errorf("Update(missing) error = %v, want ErrNotFound", err)
	}
}

func TestNewStoreCopiesSeed(t *testing.T) {
	seed := testSeed()
	store := NewStore(seed)
	status := StatusResolved

	if _, err := store.Update("1", Update{Status: &status}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	if seed[0].Status != StatusOpen {
		t.Errorf("Update mutated caller's seed slice: got %q, want %q", seed[0].Status, StatusOpen)
	}
}
