package mapping

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type xmlMapping struct{}

func (xmlMapping) Name() string {
	return "xml"
}

func (xmlMapping) Bind(req *http.Request, obj interface{}) error {
	return decodeXML(req.Body, obj)
}

func (xmlMapping) BindBody(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body), obj)
}

func decodeXML(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	return decoder.Decode(obj)
}
