package render

import (
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

type JSONP struct {
	Callback string
	Data     interface{}
}

var jsonpContentType = "application/javascript; charset=utf-8"

func Jsonp(callback string, val interface{}) Render {
	return JSONP{Data: val, Callback: callback}
}

func (r JSONP) Render(w http.ResponseWriter) (err error) {
	return WriteJSONP(w, r)
}

func (r JSONP) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonpContentType)
}

func WriteJSONP(resp http.ResponseWriter, r JSONP) error {
	writeContentType(resp, jsonContentType)
	buf, err := jsoniter.ConfigFastest.Marshal(r.Data)
	if err != nil {
		return err
	}
	io.WriteString(resp, r.Callback)
	resp.Write([]byte{'('})
	resp.Write(buf)
	_, err = resp.Write([]byte{')'})
	return err
}
