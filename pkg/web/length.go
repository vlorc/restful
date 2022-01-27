package web

import (
	"bytes"
	"io"
	"strings"
)

func Length(v interface{}) int64 {
	switch r := v.(type) {
	case int:
		return int64(r)
	case int64:
		return r
	case *bytes.Reader:
		return int64(r.Len())
	case *bytes.Buffer:
		return int64(r.Len())
	case *strings.Reader:
		return r.Size()
	case io.Seeker:
		end, _ := r.Seek(0, io.SeekEnd)
		begin, _ := r.Seek(0, io.SeekStart)
		return end - begin
	case []byte:
		return int64(len(r))
	case string:
		return int64(len(r))
	default:
		return noWritten
	}
}
