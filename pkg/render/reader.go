package render

import (
	"github.com/vlorc/restful/pkg/web"
	"io"
	"net/http"
	"strconv"
)

type Reader struct {
	ContentType   string
	ContentLength int64
	Reader        io.Reader
	Headers       map[string]string
}

func (r Reader) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	if r.ContentLength >= 0 {
		if r.Headers == nil {
			writeContentLength(w, r.ContentLength)
		} else {
			r.Headers[web.HTTP_HEADER_CONTENT_LENGTH] = strconv.FormatInt(r.ContentLength, 10)
		}
	}
	r.writeHeaders(w, r.Headers)
	_, err = io.Copy(w, r.Reader)
	return
}

func (r Reader) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, r.ContentType)
}

func (r Reader) writeHeaders(w http.ResponseWriter, headers map[string]string) {
	h := w.Header()
	for k, v := range headers {
		if h.Get(k) == "" {
			h.Set(k, v)
		}
	}
}
