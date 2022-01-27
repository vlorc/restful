package web

import (
	"regexp"
)

func Match(s ...string) func(string) bool {
	if len(s) == 0 {
		return func(string) bool {
			return true
		}
	}
	if len(s) == 1 {
		v := regexp.MustCompile(s[0])
		return func(r string) bool {
			return v.MatchString(r)
		}
	}

	v := make([]*regexp.Regexp, len(s))
	for i := range s {
		v[i] = regexp.MustCompile(s[i])
	}
	return func(r string) bool {
		for i := range v {
			if v[i].MatchString(r) {
				return true
			}
		}
		return false
	}
}
