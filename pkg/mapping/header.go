package mapping

import (
	"net/http"
)

type headerMapping struct{}

func (headerMapping) Name() string {
	return "header"
}

func (headerMapping) Bind(req *http.Request, obj interface{}) error {
	return mapHeader(obj, req.Header)
}

func mapHeader(ptr interface{}, data map[string][]string) error {
	return mapByTag(ptr, data, "header")
}
