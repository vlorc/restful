package wrap

import (
	"fmt"
	"io"
	"net/http"
	"github.com/vlorc/restful/pkg/advice"
	"github.com/vlorc/restful/pkg/binding"
	"github.com/vlorc/restful/pkg/web"
	"reflect"
	"strconv"
)

var errorType = reflect.TypeOf((*error)(nil)).Elem()

func Wrap(handlers ...interface{}) http.HandlerFunc {
	if len(handlers) == 1 {
		return wrap(handlers[0])
	}

	h := make([]http.HandlerFunc, len(handlers))

	for i := range h {
		h[i] = wrap(handlers[i])
	}

	return func(resp http.ResponseWriter, req *http.Request) {
		for i := range h {
			h[i].ServeHTTP(resp, req)
		}
	}
}

func wrap(handler interface{}) http.HandlerFunc {
	switch v := handler.(type) {
	case http.HandlerFunc:
		return v
	case func(http.ResponseWriter, *http.Request):
		return v
	case http.Handler:
		return v.ServeHTTP
	case func(io.Reader) io.ReadCloser:
		return wrapStream(v)
	case func() io.ReadCloser:
		return wrapStream(func(io.Reader) io.ReadCloser {
			return v()
		})
	case func():
		return func(resp http.ResponseWriter, req *http.Request) {
			v()
		}
	case web.HTML:
		return wrapData(v, web.HTTP_MIME_HTML_UTF8)
	case []byte:
		return wrapBytes(v)
	case web.DATA:
		return wrapBytes(v)
	case web.TEXT:
		return wrapString(v)
	case string:
		return wrapString([]byte(v))
	case web.STRING:
		return wrapString([]byte(v))
	case io.Reader:
		return wrapReader(v, web.HTTP_MIME_OCTET_STREAM)
	case web.Location:
		return func(resp http.ResponseWriter, req *http.Request) {
			http.Redirect(resp, req, string(v), http.StatusFound)
		}
	case web.Redirect:
		return func(resp http.ResponseWriter, req *http.Request) {
			http.Redirect(resp, req, v.Location, v.Code)
		}
	}

	if v := reflect.ValueOf(handler); reflect.Func == v.Kind() {
		return wrapFunc(handler)
	}

	panic(fmt.Sprintf("Unsupported handler type: %T", handler))
}

func wrapStream(f func(io.Reader) io.ReadCloser) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		if r := f(req.Body); nil != r {
			resp.WriteHeader(200)
			io.Copy(resp, r)
		}
	}
}

func wrapFunc(f interface{}) http.HandlerFunc {
	call := wrapCall(f)

	t := reflect.TypeOf(f)
	if t.NumOut() == 0 {
		return func(resp http.ResponseWriter, req *http.Request) {
			call(resp, req)
		}
	}

	o := make([]advice.Advice, t.NumOut())
	for i := range o {
		o[i] = advice.New(t.Out(i))
	}

	if t.Out(len(o)-1) != errorType {
		return func(resp http.ResponseWriter, req *http.Request) {
			r := call(resp, req)
			for i := range o {
				o[i](resp, req, r[i].Interface())
			}
		}
	}

	return func(resp http.ResponseWriter, req *http.Request) {
		r := call(resp, req)
		if !r[len(r)-1].IsNil() {
			o[len(o)-1](resp, req, r[len(r)-1].Interface())
			return
		}
		for i := range o[:len(o)-1] {
			o[i](resp, req, r[i].Interface())
		}
	}
}

func wrapCall(f interface{}) func(http.ResponseWriter, *http.Request) []reflect.Value {
	v := reflect.ValueOf(f)
	t := v.Type()
	b := make([]binding.Binding, t.NumIn())
	for i := range b {
		b[i] = binding.New(t.In(i))
	}

	return func(resp http.ResponseWriter, req *http.Request) []reflect.Value {
		p := make([]reflect.Value, len(b))
		t := v.Type()
		for i := range p {
			p[i] = b[i](resp, req, t.In(i))
		}
		return v.Call(p)
	}
}

func wrapData(data []byte, typ string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		h := resp.Header()
		h[web.HTTP_HEADER_CONTENT_LENGTH] = []string{strconv.Itoa(len(data))}
		h[web.HTTP_HEADER_CONTENT_TYPE] = []string{typ}
		resp.WriteHeader(http.StatusOK)
		resp.Write(data)
	}
}

func wrapBytes(data []byte) http.HandlerFunc {
	return wrapData(data, web.HTTP_MIME_OCTET_STREAM)
}

func wrapString(data []byte) http.HandlerFunc {
	return wrapData(data, web.HTTP_MIME_TEXT_UTF8)
}

func wrapReader(r io.Reader, typ string) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		h := resp.Header()
		if l := web.Length(r); l >= 0 {
			h[web.HTTP_HEADER_CONTENT_LENGTH] = []string{strconv.FormatInt(l, 10)}
		}
		h[web.HTTP_HEADER_CONTENT_TYPE] = []string{typ}
		resp.WriteHeader(http.StatusOK)
		io.Copy(resp, r)
	}
}
