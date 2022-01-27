package render

import (
	"github.com/json-iterator/go"
	"net/http"
)

type JSON struct {
	Data interface{}
}

var jsonContentType = "application/json; charset=utf-8"

func Json(val interface{}) Render {
	return JSON{Data: val}
}

func (r JSON) Render(w http.ResponseWriter) (err error) {
	return WriteJSON(w, r.Data)
}

func (r JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func WriteJSON(resp http.ResponseWriter, obj interface{}) error {
	writeContentType(resp, jsonContentType)
	buf, err := jsoniter.ConfigFastest.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = resp.Write(buf)
	return err
}
