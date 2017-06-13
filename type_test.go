package pbast

import (
	"testing"
)

func TestIsExactlySameType(t *testing.T) {
	table := []struct {
		t1       Type
		t2       Type
		expected bool
	}{
		// same name and same type
		{
			t1:       &Message{Name: "M1"},
			t2:       &Message{Name: "M1"},
			expected: true,
		},
		// same name but different type
		{
			t1:       &Message{Name: "M1"},
			t2:       &Enum{Name: "M1"},
			expected: false,
		},
		{
			t1: &Message{
				Name: "M1",
				Fields: []*MessageField{
					&MessageField{Type: "F1", Name: "field", Index: 1},
				},
			},
			t2:       &Message{Name: "M1"},
			expected: false,
		},
	}

	for x, d := range table {
		if actual := IsSameType(d.t1, d.t2); actual != d.expected {
			t.Errorf("#%d: got %t, want %t", x, actual, d.expected)
		}
	}
}
