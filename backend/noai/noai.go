package main

import (
	"fmt"
	"slices"
)

type Conversation struct {
	ID       string
	Customer string
	Priority string
	Status   string
}

// NO-AI TASK
func SummarizeConversations(
	conversations []Conversation,
) map[string]int {
	var status_priority []string

	for _, user := range conversations {
		if len(user.Status) > 0 && len(user.Priority) > 0 {
			status_priority = append(status_priority, fmt.Sprintf("%s_%s", user.Status, user.Priority))
		}
	}
	combined := make(map[string]int)
	copied := slices.Clone(status_priority)

	slices.Sort(copied)

	uniqueStrings := slices.Compact(copied)
	for _, val := range uniqueStrings {
		combined[val] = countOccurrences(status_priority, val)
	}

	return combined
}

func countOccurrences(slice []string, target string) int {
	count := 0
	for _, v := range slice {
		if v == target {
			count++
		}
	}
	return count
}

func main() {
	slice := []Conversation{{ID: "conv-001", Customer: "John Carter", Status: "IN_PROGRESS", Priority: "LOW"}, {ID: "conv-002", Customer: "Maria Lopez", Status: "IN_PROGRESS", Priority: "MEDIUM"}, {ID: "conv-003", Customer: "John Nguyen", Status: "IN_PROGRESS", Priority: "MEDIUM"}, {ID: "conv-004", Customer: "Priya Patel", Status: "IN_PROGRESS", Priority: "MEDIUM"}, {ID: "conv-005", Customer: "Tom Becker", Status: "RESOLVED", Priority: "LOW"}, {ID: "conv-006", Customer: "Aisha Khan", Status: "IN_PROGRESS", Priority: "HIGH"}, {ID: "conv-007", Customer: "Lucas Moreau", Status: "RESOLVED", Priority: "MEDIUM"}, {ID: "conv-008", Customer: "Emily Watson", Status: "RESOLVED", Priority: "LOW"}}

	result := SummarizeConversations(slice)

	fmt.Print(result)
}
