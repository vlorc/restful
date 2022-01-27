package router

import "strings"

func Pattern(prefix string) func(string2 string) string {
	return func(pattern string) string {
		p := prefix + pattern
		if !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		if p == "/" {
			return p
		}
		return strings.TrimSuffix(p, "/")
	}
}
