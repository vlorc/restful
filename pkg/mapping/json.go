package mapping

import (
	"bytes"
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"net/http"
)

type jsonMapping struct{}

func (jsonMapping) Name() string {
	return "json"
}

func (jsonMapping) Bind(req *http.Request, obj interface{}) error {
	if req == nil || req.Body == nil {
		return fmt.Errorf("invalid request")
	}
	return decodeJSON(req.Body, obj)
}

func (jsonMapping) BindBody(body []byte, obj interface{}) error {
	return decodeJSON(bytes.NewReader(body), obj)
}

func decodeJSON(r io.Reader, obj interface{}) error {
	return jsoniter.ConfigFastest.NewDecoder(r).Decode(obj)
}
