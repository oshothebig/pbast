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
		{"ETHERNET", []string{"ETHERNET"}},
		{"OCS", []string{"OCS"}},
		{"aToZ", []string{"a", "To", "Z"}},
		{"100G", []string{"100", "G"}},
		{"Ipv4", []string{"Ipv4"}},
		{"IPV4", []string{"IP", "V4"}},
		// looks strange, but it's because we don't have a dictonary for abbreviations
		{"IPv4", []string{"I", "Pv4"}},
	}

	for x, d := range table {
		if actual := guessElements(d.in); !reflect.DeepEqual(actual, d.expected) {
			t.Errorf("#%d: got %v, want %v", x, actual, d.expected)
		}
	}
}
