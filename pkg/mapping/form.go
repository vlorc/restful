package mapping

import (
	"net/http"
)

const defaultMemory = 32 << 20

type formMapping struct{}
type formPostMapping struct{}
type formMultipartMapping struct{}

func (formMapping) Name() string {
	return "form"
}

func (formMapping) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		if err != http.ErrNotMultipart {
			return err
		}
	}
	return mapForm(obj, req.Form)
}

func (formPostMapping) Name() string {
	return "form-urlencoded"
}

func (formPostMapping) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	return mapForm(obj, req.PostForm)
}

func (formMultipartMapping) Name() string {
	return "multipart/form-data"
}

func (formMultipartMapping) Bind(req *http.Request, obj interface{}) error {
	if err := req.ParseMultipartForm(defaultMemory); err != nil {
		return err
	}

	return mapForm(obj, req.PostForm)
}

func mapForm(ptr interface{}, data map[string][]string) error {
	return mapByTag(ptr, data, "form")
}
