package yang

import (
	"errors"
	"net/url"
	"strings"
	"unicode"
)

func guessElements(s string) []string {
	// URL based
	if strings.Contains(s, "://") {
		if u, err := url.Parse(s); err == nil {
			elems := strings.Split(u.Hostname(), ".")

			// reverse
			for i := 0; i < len(elems)/2; i++ {
				opposite := len(elems) - 1 - i
				elems[i], elems[opposite] = elems[opposite], elems[i]
			}

			path := strings.Trim(u.Path, "/")
			if len(path) != 0 {
				elems = append(elems, strings.Split(path, "/")...)
			}
			return elems
		}
	}

	// separator based
	if sep, err := guessSeparator(s, []string{".", ":", "_", "-"}); err == nil {
		return strings.Split(s, sep)
	}

	// CamelCase based
	if segment, err := splitCamelCase(s); err == nil {
		return segment
	}

	// give up guessing
	return []string{s}
}

func guessSeparator(s string, seps []string) (string, error) {
	found := []string{}
	for _, sep := range seps {
		if strings.Contains(s, sep) {
			found = append(found, sep)
		}
	}

	if len(found) == 0 {
		return "", errors.New("not found")
	}
	if len(found) > 1 {
		return "", errors.New("ambiguous")
	}

	return found[0], nil
}

func splitCamelCase(s string) ([]string, error) {
	indexes := upperCaseIndexes(s)

	// no upper case character found
	if len(indexes) == 0 {
		return []string{s}, nil
	}

	indexes = complementFirstCharacter(indexes)
	indexes = complementLastCharacter(indexes, s)
	indexes = handleAbbreviation(indexes)
	return splitByIndexes(s, indexes), nil
}

func upperCaseIndexes(s string) []int {
	var indexes []int
	for i, ch := range s {
		if unicode.IsUpper(ch) {
			indexes = append(indexes, i)
		}
	}

	return indexes
}

func complementFirstCharacter(indexes []int) []int {
	if indexes[0] != 0 {
		return append([]int{0}, indexes...)
	}

	return indexes
}

func complementLastCharacter(indexes []int, s string) []int {
	if indexes[len(indexes)-1] == len(s) {
		return indexes
	}

	return append(indexes, len(s))
}

func handleAbbreviation(indexes []int) []int {
	filtered := []int{0}
	for pos := 1; pos < len(indexes)-1; pos++ {
		if indexes[pos+1]-indexes[pos] == 1 {
			continue
		}

		filtered = append(filtered, indexes[pos])
	}
	filtered = append(filtered, indexes[len(indexes)-1])

	return filtered
}

func splitByIndexes(s string, indexes []int) []string {
	if len(indexes) == 0 {
		return []string{s}
	}

	var splitted []string
	for pos := 1; pos < len(indexes); pos++ {
		start := indexes[pos-1]
		end := indexes[pos]
		splitted = append(splitted, s[start:end])
	}

	return splitted
}

func CamelCase(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.Title(e)
	}

	return strings.Join(elems, "")
}

func underscoreCase(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.ToLower(e)
	}

	return strings.Join(elems, "_")
}

func constantName(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.ToUpper(e)
	}

	return strings.Join(elems, "_")
}

func buildName(base, common, suffix string) string {
	lowerBase := strings.ToLower(base)
	lowerCommon := strings.ToLower(common)
	if strings.HasSuffix(lowerBase, lowerCommon) {
		return base + suffix
	}

	return base + common + suffix
}

type stringSet map[string]struct{}

func newStringSet() stringSet {
	return map[string]struct{}{}
}

func newStringSetWith(ss []string) stringSet {
	set := newStringSet()
	for _, s := range ss {
		set[s] = struct{}{}
	}
	return set
}

func (s stringSet) contains(element string) bool {
	_, ok := s[element]
	return ok
}

func (s stringSet) add(element string) {
	s[element] = struct{}{}
}

func (s stringSet) remove(element string) {
	delete(s, element)
}
