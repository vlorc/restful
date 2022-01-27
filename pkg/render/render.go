package render

import (
	"net/http"
	"github.com/vlorc/restful/pkg/web"
	"strconv"
)

type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(http.ResponseWriter)
}

func writeContentLength(w http.ResponseWriter, v int64) {
	if h, k := w.Header(), web.HTTP_HEADER_CONTENT_LENGTH; len(h[k]) == 0 {
		h[k] = []string{strconv.FormatInt(v, 10)}
	}
}

func writeContentType(w http.ResponseWriter, v string) {
	if h, k := w.Header(), web.HTTP_HEADER_CONTENT_TYPE; len(h[k]) == 0 {
		h[k] = []string{v}
	}
}

func writeNoCache(w http.ResponseWriter, v string) {
	if h, k := w.Header(), web.HTTP_HEADER_CACHE_CONTROL; len(h[k]) == 0 {
		h[k] = []string{v}
	}
}
