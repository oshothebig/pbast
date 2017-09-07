package yang

import "testing"
import "reflect"

func TestGuessElements(t *testing.T) {
	table := []struct {
		in       string
		expected []string
	}{
		{"a.b.c", []string{"a", "b", "c"}},
		{"a:b:c", []string{"a", "b", "c"}},
		{"a_b_c", []string{"a", "b", "c"}},
		{"a-b-c", []string{"a", "b", "c"}},
		{"http://example/", []string{"example"}},
		{"http://www.example.com", []string{"com", "example", "www"}},
		{"http://www.example.com/", []string{"com", "example", "www"}},
		{"example.com", []string{"example", "com"}},
		{"inService", []string{"in", "Service"}},
		{"IPAddress", []string{"IP", "Address"}},
		{"HTMLRender", []string{"HTML", "Render"}},
		{"IpAddress", []string{"Ip", "Address"}},
	}

	for x, d := range table {
		if actual := guessElements(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
