package render

import (
	"encoding/xml"
	"net/http"
)

type XML struct {
	Data interface{}
}

var xmlContentType = "application/xml; charset=utf-8"

func Xml(val interface{}) Render {
	return XML{Data: val}
}

func (r XML) Render(w http.ResponseWriter) (err error) {
	return WriteJSON(w, r.Data)
}

func (r XML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, xmlContentType)
}

func WriteXML(resp http.ResponseWriter, obj interface{}) error {
	writeContentType(resp, xmlContentType)
	buf, err := xml.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = resp.Write(buf)
	return err
}
