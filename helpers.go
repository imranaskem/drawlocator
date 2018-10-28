package main

import (
	"errors"
	"regexp"
	"strings"
)

func standardisePlace(place string) (string, error) {
	place = strings.ToLower(place)

	switch {
	case strings.Contains(place, "baker"):
		return "Baker Street", nil

	case strings.Contains(place, "sick"):
		return "Sick", nil

	case strings.Contains(place, "weston"):
		return "Weston Street", nil

	case strings.Contains(place, "holiday"):
		return "Holiday", nil

	case strings.Contains(place, "client"):
		return "Client Office", nil

	case strings.Contains(place, "home"):
		return "Working from Home", nil
	}

	return place, errors.New("Invalid place")
}

func getUserID(text string) (string, error) {
	r := regexp.MustCompile("\\@([^\\|]+)\\|")

	s := r.FindStringSubmatch(text)

	return s[1], nil
}

func comparePeople(a, b []person) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
