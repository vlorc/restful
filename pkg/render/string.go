package render

import (
	"fmt"
	"github.com/vlorc/restful/pkg/web"
	"io"
	"net/http"
)

type STRING struct {
	Format string
	Data   []interface{}
}

var stringContentType = web.HTTP_HEADER_CONTENT_TYPE

func String(format string, val ...interface{}) Render {
	return STRING{Data: val, Format: format}
}

func (r STRING) Render(w http.ResponseWriter) (err error) {
	return WriteJSON(w, r.Data)
}

func (r STRING) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, stringContentType)
}

func WriteString(resp http.ResponseWriter, r STRING) error {
	writeContentType(resp, stringContentType)
	if len(r.Data) > 0 {
		_, err := fmt.Fprintf(resp, r.Format, r.Data...)
		return err
	}
	_, err := io.WriteString(resp, r.Format)
	return err
}
