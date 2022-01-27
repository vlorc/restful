package mapping

import "net/http"

type queryMapping struct{}

func (queryMapping) Name() string {
	return "query"
}

func (queryMapping) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	return mapQuery(obj, values)
}

func mapQuery(ptr interface{}, data map[string][]string) error {
	return mapByTag(ptr, data, "form")
}
