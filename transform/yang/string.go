package yang

import (
	"errors"
	"net/url"
	"strings"
	"unicode"
)

func CamelCase(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.Title(e)
	}

	return strings.Join(elems, "")
}

func snakeCase(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.ToLower(e)
	}

	return strings.Join(elems, "_")
}

func constantCase(s string) string {
	elems := guessElements(s)
	for x, e := range elems {
		elems[x] = strings.ToUpper(e)
	}

	return strings.Join(elems, "_")
}

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
	if len(s) == 0 {
		return []string{s}, nil
	}

	var ret []string
	runes := []rune(s)
	segment := []rune{runes[0]}
	for i := 1; i < len(runes); i++ {
		current := runes[i]
		previous := runes[i-1]

		if !unicode.IsUpper(current) {
			segment = append(segment, current)
			continue
		}

		if !unicode.IsUpper(previous) {
			ret = append(ret, string(segment))
			segment = []rune{current}
			continue
		}

		if i+1 == len(runes) {
			segment = append(segment, current)
			continue
		}

		next := runes[i+1]
		if unicode.IsUpper(next) {
			segment = append(segment, current)
			continue
		}

		ret = append(ret, string(segment))
		segment = []rune{current}
	}
	ret = append(ret, string(segment))

	return ret, nil
}

func buildName(base, common, suffix string) string {
	lowerBase := strings.ToLower(base)
	lowerCommon := strings.ToLower(common)
	if strings.HasSuffix(lowerBase, lowerCommon) {
		return base + suffix
	}

	return base + common + suffix
}
