package main

import (
	"reflect"
	"testing"
)

var conversations = []Conversation{{ID: "conv-001", Customer: "John Carter", Status: "", Priority: "LOW"}, {ID: "conv-002", Customer: "Maria Lopez", Status: "IN_PROGRESS", Priority: ""}, {ID: "conv-003", Customer: "John Nguyen", Status: "IN_PROGRESS", Priority: "MEDIUM"}, {ID: "conv-004", Customer: "Priya Patel", Status: "IN_PROGRESS", Priority: "MEDIUM"}, {ID: "conv-005", Customer: "Tom Becker", Status: "RESOLVED", Priority: "LOW"}, {ID: "conv-006", Customer: "Aisha Khan", Status: "IN_PROGRESS", Priority: "HIGH"}, {ID: "conv-007", Customer: "Lucas Moreau", Status: "RESOLVED", Priority: "MEDIUM"}, {ID: "conv-008", Customer: "Emily Watson", Status: "RESOLVED", Priority: "LOW"}}

func TestEmptiness(t *testing.T) {
	got := SummarizeConversations(conversations)
	want := map[string]int{
		"IN_PROGRESS_HIGH": 1, "IN_PROGRESS_MEDIUM": 2, "RESOLVED_LOW": 2, "RESOLVED_MEDIUM": 1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestIdempotency(t *testing.T) {
	gotFirst := SummarizeConversations(conversations)
	gotSecond := SummarizeConversations(conversations)

	if !reflect.DeepEqual(gotFirst, gotSecond) {
		t.Errorf("got %v, want %v", gotFirst, gotSecond)
	}
}
