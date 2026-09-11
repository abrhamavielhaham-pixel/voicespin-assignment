package conversation

import "time"

func Seed() []Conversation {
	return []Conversation{
		{
			ID:            "conv-001",
			CustomerName:  "John Carter",
			CustomerEmail: "john.carter@example.com",
			Subject:       "Cannot reset my password",
			Status:        StatusOpen,
			Priority:      PriorityHigh,
			CreatedAt:     time.Date(2026, 9, 1, 9, 15, 0, 0, time.UTC),
		},
		{
			ID:            "conv-002",
			CustomerName:  "Maria Lopez",
			CustomerEmail: "maria.lopez@example.com",
			Subject:       "Invoice charged twice",
			Status:        StatusOpen,
			Priority:      PriorityHigh,
			CreatedAt:     time.Date(2026, 9, 2, 11, 40, 0, 0, time.UTC),
		},
		{
			ID:            "conv-003",
			CustomerName:  "John Nguyen",
			CustomerEmail: "john.nguyen@acme.io",
			Subject:       "Question about API rate limits",
			Status:        StatusInProgress,
			Priority:      PriorityMedium,
			CreatedAt:     time.Date(2026, 9, 3, 8, 5, 0, 0, time.UTC),
		},
		{
			ID:            "conv-004",
			CustomerName:  "Priya Patel",
			CustomerEmail: "priya.patel@example.com",
			Subject:       "Export to CSV fails for large reports",
			Status:        StatusOpen,
			Priority:      PriorityMedium,
			CreatedAt:     time.Date(2026, 9, 4, 14, 25, 0, 0, time.UTC),
		},
		{
			ID:            "conv-005",
			CustomerName:  "Tom Becker",
			CustomerEmail: "tom.becker@contoso.com",
			Subject:       "Feature request: dark mode",
			Status:        StatusResolved,
			Priority:      PriorityLow,
			CreatedAt:     time.Date(2026, 9, 5, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:            "conv-006",
			CustomerName:  "Aisha Khan",
			CustomerEmail: "aisha.khan@example.com",
			Subject:       "Login fails after SSO migration",
			Status:        StatusInProgress,
			Priority:      PriorityHigh,
			CreatedAt:     time.Date(2026, 9, 6, 16, 45, 0, 0, time.UTC),
		},
		{
			ID:            "conv-007",
			CustomerName:  "Lucas Moreau",
			CustomerEmail: "lucas.moreau@example.com",
			Subject:       "Wrong currency shown on dashboard",
			Status:        StatusResolved,
			Priority:      PriorityMedium,
			CreatedAt:     time.Date(2026, 9, 7, 9, 30, 0, 0, time.UTC),
		},
		{
			ID:            "conv-008",
			CustomerName:  "Emily Watson",
			CustomerEmail: "emily.watson@example.com",
			Subject:       "How do I invite a teammate?",
			Status:        StatusResolved,
			Priority:      PriorityLow,
			CreatedAt:     time.Date(2026, 9, 8, 13, 10, 0, 0, time.UTC),
		},
	}
}
