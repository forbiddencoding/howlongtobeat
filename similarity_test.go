package howlongtobeat

import "testing"

func Test_calculateJaccardSimilarity(t *testing.T) {
	tables := []struct {
		gameTitle1 string
		gameTitle2 string
		expected   float64
	}{
		{"Super Mario", "Super Mario", 1},
		{"Super Mario", "Mario Super", 1},
		{"Super Mario Bros", "Super Luigi Bros", 0.5},
		{"Super Mario Bros", "Super Luigi BROS", 0.5}, // case insensitive
		{"Halo", "Call Of Duty", 0},
		{"", "", 0},
		{"Final Fantasy", "Final Fast", 0.33},
		{"Multi word title here", "Multi multi word word title title here here", 1.0}, // repeating words
	}

	for _, table := range tables {
		result := calculateJaccardSimilarity(table.gameTitle1, table.gameTitle2)
		if result != table.expected {
			t.Errorf("calculateJaccardSimilarity(%s, %s) was incorrect, got: %f, want: %f.", table.gameTitle1, table.gameTitle2, result, table.expected)
		}
	}
}
