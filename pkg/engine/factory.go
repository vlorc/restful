package engine

import (
	"github.com/vlorc/restful/pkg/router"
	"net/http"
)

const DETAULT_NAME = "default"

var mapping = map[string]func(...func(http.Handler) http.Handler) router.Router{}

func Register(name string, build func(...func(http.Handler) http.Handler) router.Router) {
	mapping[name] = build
}

func New(name string, middles ...func(http.Handler) http.Handler) router.Router {
	if b, ok := mapping[name]; ok {
		return b(middles...)
	}
	return nil
}

func Default(middles ...func(http.Handler) http.Handler) router.Router {
	if b, ok := mapping[DETAULT_NAME]; ok {
		return b(middles...)
	}
	return nil
}
