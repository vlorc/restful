package advice

import (
	"context"
	"io"
	"net/http"
	"github.com/vlorc/restful/pkg/render"
	"github.com/vlorc/restful/pkg/web"
	"reflect"
	"strconv"
)

var renderType = reflect.TypeOf((*render.Render)(nil)).Elem()

func types() map[reflect.Type][]Advice {
	m := map[reflect.Type][]Advice{
		reflect.TypeOf((*web.Code)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.Code); ok && 0 != v {
				resp.WriteHeader(int(v))
			}
			return val
		}},
		reflect.TypeOf((*func() web.Code)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() web.Code); ok && nil != v {
				if r := v(); 0 != r {
					resp.WriteHeader(int(r))
				}
			}
			return val
		}},
		reflect.TypeOf((*int)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(int); ok && 0 != v {
				resp.WriteHeader(v)
			}
			return val
		}},
		reflect.TypeOf((*func() int)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() int); ok && nil != v {
				if r := v(); 0 != r {
					resp.WriteHeader(r)
				}
			}
			return val
		}},
		reflect.TypeOf((*int64)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(int64); ok && v >= 0 {
				resp.Header().Set(web.HTTP_HEADER_CONTENT_LENGTH, strconv.FormatInt(v, 10))
			}
			return val
		}},
		reflect.TypeOf((*func() int64)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() int64); ok && nil != v {
				if r := v(); r >= 0 {
					resp.Header().Set(web.HTTP_HEADER_CONTENT_LENGTH, strconv.FormatInt(r, 10))
				}
			}
			return val
		}},
		reflect.TypeOf((*io.Reader)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(io.Reader); ok && nil != v {
				writeContentLength(resp, v)
				io.Copy(resp, v)
			}
			return val
		}},
		reflect.TypeOf((*io.ReadCloser)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(io.ReadCloser); ok && nil != v {
				defer v.Close()
				writeContentLength(resp, v)
				io.Copy(resp, v)
			}
			return val
		}},
		reflect.TypeOf((*func() io.Reader)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() io.Reader); ok && nil != v {
				if r := v(); nil != r {
					writeContentLength(resp, r)
					io.Copy(resp, r)
				}
			}
			return val
		}},
		reflect.TypeOf((*func() io.ReadCloser)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() io.ReadCloser); ok && nil != v {
				if r := v(); nil != r {
					defer r.Close()
					writeContentLength(resp, r)
					io.Copy(resp, r)
				}
			}
			return val
		}},
		reflect.TypeOf((*[]byte)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.([]byte); ok && nil != v {
				writeContentLength(resp, v)
				resp.Write(v)
			}
			return val
		}},
		reflect.TypeOf((*string)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(string); ok && "" != v {
				writeContentLength(resp, v)
				io.WriteString(resp, v)
			}
			return val
		}},
		reflect.TypeOf((*func() []byte)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() []byte); ok && nil != v {
				if r := v(); nil != r {
					writeContentLength(resp, r)
					resp.Write(r)
				}
			}
			return val
		}},
		reflect.TypeOf((*func() string)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() string); ok && nil != v {
				if r := v(); "" != r {
					writeContentLength(resp, r)
					io.WriteString(resp, r)
				}
			}
			return val
		}},
		reflect.TypeOf((*web.ContentLength)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.ContentLength); ok && v >= 0 {
				writeContentLength(resp, int64(v))
			}
			return val
		}},
		reflect.TypeOf((*func() web.ContentLength)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() web.ContentLength); ok && nil != v {
				if r := v(); r >= 0 {
					writeContentLength(resp, int64(r))
				}
			}
			return val
		}},
		reflect.TypeOf((*web.ContentType)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.ContentType); ok && "" != v {
				writeContentType(resp, string(v))
			}
			return val
		}},
		reflect.TypeOf((*func() web.ContentType)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() web.ContentType); ok && nil != v {
				if r := v(); "" != r {
					writeContentType(resp, string(r))
				}
			}
			return val
		}},
		reflect.TypeOf((*web.HTML)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.HTML); ok && nil != v {
				writeData(resp, http.StatusOK, web.HTTP_MIME_HTML_UTF8, v)
			}
			return val
		}},
		reflect.TypeOf((*web.TEXT)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.TEXT); ok && nil != v {
				writeData(resp, http.StatusOK, web.HTTP_MIME_TEXT_UTF8, v)
			}
			return val
		}},
		reflect.TypeOf((*web.STRING)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.STRING); ok && "" != v {
				writeString(resp, http.StatusOK, string(v))
			}
			return val
		}},
		reflect.TypeOf((*web.DATA)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.DATA); ok && nil != v {
				writeData(resp, http.StatusOK, web.HTTP_MIME_OCTET_STREAM, v)
			}
			return val
		}},
		reflect.TypeOf((*web.Redirect)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.Redirect); ok && 0 != v.Code && "" != v.Location {
				http.Redirect(resp, req, v.Location, v.Code)
			}
			return val
		}},
		reflect.TypeOf((*web.Redirect)(nil)): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(*web.Redirect); ok && nil != v && 0 != v.Code && "" != v.Location {
				http.Redirect(resp, req, v.Location, v.Code)
			}
			return val
		}},
		reflect.TypeOf((*web.Location)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(web.Location); ok && "" != v {
				http.Redirect(resp, req, string(v), http.StatusFound)
			}
			return val
		}},
		reflect.TypeOf((*func(io.Writer) bool)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func(io.Writer) bool); ok && nil != v {
				bindStream(resp, req, v)
			}
			return val
		}},

		reflect.TypeOf((*func(ctx context.Context) io.Reader)(nil)).Elem(): {bindStreamReader},
		reflect.TypeOf((*chan io.Reader)(nil)).Elem():                      {bindChanReader},
		renderType: {bindRender},
		reflect.TypeOf((*func() render.Render)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() render.Render); ok && nil != v {
				bindRender(resp, req, v())
			}
			return val
		}},
		reflect.TypeOf((*chan render.Render)(nil)).Elem():                  {bindChanRender},
		reflect.TypeOf((*func(context.Context) render.Render)(nil)).Elem(): {bindStreamRender},
		reflect.TypeOf((*func())(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func()); ok && nil != v {
				v()
			}
			return val
		}},
		errorType: {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(error); ok && nil != v {
				bindError(resp, req, v)
			}
			return val
		}},
		reflect.TypeOf((*func() error)(nil)).Elem(): {func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if v, ok := val.(func() error); ok && nil != v {
				bindError(resp, req, v())
			}
			return val
		}},
	}

	return m
}

func bindStreamReader(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	v, _ := val.(func(context.Context) io.Reader)
	if nil == v {
		return val
	}

	if r, ok := resp.(web.ResponseWriter); ok && !bodyAllowedForStatus(r.Status()) {
		r.WriteHeaderNow()
		return val
	}

	f := resp.(http.Flusher)
	ctx := req.Context()
	for {
		r := v(ctx)
		if nil == r {
			break
		}

		_, err := io.Copy(resp, r)
		f.Flush()
		if nil != err {
			break
		}
	}

	return val
}

func bindChanReader(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	v, _ := val.(chan io.Reader)
	if nil == v {
		return val
	}

	if r, ok := resp.(web.ResponseWriter); ok && !bodyAllowedForStatus(r.Status()) {
		r.WriteHeaderNow()
		return val
	}

	f := resp.(http.Flusher)
	ctx := req.Context()
	for {
		select {
		case <-ctx.Done():
			break
		case r, ok := <-v:
			if !ok && nil == r {
				break
			}
			_, err := io.Copy(resp, r)
			f.Flush()
			if nil != err {
				break
			}
		}
	}

	return val
}

func Render(t reflect.Type) Advice {
	if nil == t || t == renderType || t.ConvertibleTo(renderType) {
		return bindRender
	}
	return nil
}

func bindStreamRender(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	if r, ok := resp.(web.ResponseWriter); ok && !bodyAllowedForStatus(r.Status()) {
		r.WriteHeaderNow()
		return val
	}

	v, ok := val.(func(context.Context) render.Render)
	if !ok || nil == v {
		return val
	}

	f := resp.(http.Flusher)
	ctx := req.Context()
	for {
		r := v(ctx)
		if nil == r {
			break
		}

		err := r.Render(resp)
		f.Flush()
		if nil != err {
			break
		}
	}

	return val
}

func bindChanRender(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	v, _ := val.(chan render.Render)
	if nil == v {
		return val
	}

	if r, ok := resp.(web.ResponseWriter); ok && !bodyAllowedForStatus(r.Status()) {
		r.WriteHeaderNow()
		return val
	}

	f := resp.(http.Flusher)
	ctx := req.Context()
	for {
		select {
		case <-ctx.Done():
			break
		case r, ok := <-v:
			if !ok && nil == r {
				break
			}
			err := r.Render(resp)
			f.Flush()
			if nil != err {
				break
			}
		}
	}

	return val
}

func bindRender(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	if v, ok := val.(render.Render); ok && nil != v {
		if r, ok := resp.(web.ResponseWriter); ok && !bodyAllowedForStatus(r.Status()) {
			v.WriteContentType(r)
			r.WriteHeaderNow()
		} else {
			v.Render(resp)
		}
	}
	return val
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

func writeData(w http.ResponseWriter, c int, t string, v []byte) {
	writeContentType(w, t)
	writeContentLength(w, len(v))
	w.WriteHeader(c)
	w.Write(v)
}

func writeString(w http.ResponseWriter, c int, v string) {
	writeContentType(w, web.HTTP_MIME_TEXT_UTF8)
	writeContentLength(w, len(v))
	w.WriteHeader(c)
	io.WriteString(w, v)
}

func writeContentLength(w http.ResponseWriter, v interface{}) {
	if h, k := w.Header(), web.HTTP_HEADER_CONTENT_LENGTH; len(h[k]) == 0 {
		if r := web.Length(v); r >= 0 {
			h[k] = []string{strconv.FormatInt(r, 10)}
		}
	}
}

func writeContentType(w http.ResponseWriter, v string) {
	if h, k := w.Header(), web.HTTP_HEADER_CONTENT_TYPE; len(h[k]) == 0 {
		h[k] = []string{v}
	}
}

func bindStream(resp http.ResponseWriter, req *http.Request, step func(w io.Writer) bool) bool {
	f := resp.(http.Flusher)
	for {
		select {
		case <-req.Context().Done():
			return true
		default:
			keep := step(resp)
			f.Flush()
			if !keep {
				return false
			}
		}
	}
}

func Reverse(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	if v, ok := val.(func() interface{}); ok && nil != v {
		Call(resp, req, v())
	}
	return val
}

func bindError(resp http.ResponseWriter, req *http.Request, err error) {
	if nil != err {
		writeString(resp, http.StatusInternalServerError, err.Error())
	}
}

func Panic(_ http.ResponseWriter, _ *http.Request, val interface{}) interface{} {
	panic(val)
}

func Error() Advice {
	return New(errorType)
}

func Abort(http.ResponseWriter, *http.Request, interface{}) interface{} {
	panic(http.ErrAbortHandler)
}
