package yang

import (
	"errors"
	"net/url"
	"strings"
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
