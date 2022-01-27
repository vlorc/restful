package binding

import (
	"bytes"
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"github.com/vlorc/restful/pkg/web"
	"reflect"
	"strconv"
)

func types() map[reflect.Type]Binding {
	return map[reflect.Type]Binding{
		reflect.TypeOf((*http.Request)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req)
		},
		reflect.TypeOf((*http.Request)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(*req)
		},
		reflect.TypeOf((*http.Response)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(*req.Response)
		},
		reflect.TypeOf((*http.Response)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Response)
		},
		reflect.TypeOf((*http.ResponseWriter)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(resp)
		},
		reflect.TypeOf((*web.ResponseWriter)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			if v, ok := resp.(web.ResponseWriter); ok && nil != v {
				return reflect.ValueOf(v)
			}
			return reflect.ValueOf(web.NewResponseWriter(resp))
		},
		reflect.TypeOf((*http.Flusher)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			f, _ := resp.(http.Flusher)
			return reflect.ValueOf(f)
		},
		reflect.TypeOf((*http.Hijacker)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			h, _ := resp.(http.Hijacker)
			return reflect.ValueOf(h)
		},
		reflect.TypeOf((*io.Reader)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Body.(io.Reader))
		},
		reflect.TypeOf((*io.ReadCloser)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Body.(io.ReadCloser))
		},
		reflect.TypeOf((*func() (io.ReadSeeker, error))(nil)).Elem(): bindReadSeeker,
		reflect.TypeOf((*[]byte)(nil)).Elem():                        bindBytes,
		reflect.TypeOf((*io.Writer)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(resp.(io.Writer))
		},
		reflect.TypeOf((*io.WriteCloser)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Body.(io.WriteCloser))
		},
		reflect.TypeOf((*url.URL)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.URL)
		},
		reflect.TypeOf((*url.URL)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(*req.URL)
		},
		reflect.TypeOf((http.Header)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Header)
		},
		reflect.TypeOf((url.Values)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Form)
		},
		reflect.TypeOf((*multipart.Form)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.MultipartForm)
		},
		reflect.TypeOf((*web.RemoteAddr)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.RemoteAddr(req.RemoteAddr))
		},
		reflect.TypeOf((*web.RequestURI)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.RequestURI(req.RequestURI))
		},
		reflect.TypeOf((*tls.ConnectionState)(nil)): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.TLS)
		},
		reflect.TypeOf((*context.Context)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Context())
		},
		reflect.TypeOf((*web.ContentLength)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			v := req.ContentLength
			if 0 == v {
				v, _ = strconv.ParseInt(req.Header.Get(web.HTTP_HEADER_CONTENT_LENGTH), 10, 0)
			}
			return reflect.ValueOf(web.ContentLength(v))
		},
		reflect.TypeOf((*web.ContentType)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.ContentType(req.Header.Get(web.HTTP_HEADER_CONTENT_TYPE)))
		},
		reflect.TypeOf((*web.TransferEncoding)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.TransferEncoding(req.TransferEncoding))
		},
		reflect.TypeOf((*web.UserAgent)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.UserAgent(req.UserAgent()))
		},
		reflect.TypeOf((*web.RawQuery)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.RawQuery(req.URL.RawQuery))
		},
		reflect.TypeOf((*web.RawPath)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.RawPath(req.URL.RawPath))
		},
		reflect.TypeOf((*web.Path)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.Path(req.URL.Path))
		},
		reflect.TypeOf((*web.Query)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			q := req.URL.Query()
			return reflect.ValueOf(web.Query(func(k string) string {
				return q.Get(k)
			}))
		},
		reflect.TypeOf((*web.Header)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.Header(func(k string) string {
				return req.Header.Get(k)
			}))
		},
		reflect.TypeOf((*web.Method)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.Method(req.Method))
		},
		reflect.TypeOf((*web.Host)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.Host(req.Host))
		},
		reflect.TypeOf((*[]*http.Cookie)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(req.Cookies())
		},
		reflect.TypeOf((*func(string) *http.Cookie)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(getCookie(req))
		},
		reflect.TypeOf((*web.Cookie)(nil)).Elem(): func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			return reflect.ValueOf(web.Cookie(getCookie(req)))
		},
	}
}

func bindBytes(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
	body, err := ioutil.ReadAll(req.Body)
	if err == nil {
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
	}
	return reflect.ValueOf(body)
}

func bindReadSeeker(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
	body, err := ioutil.ReadAll(req.Body)
	r := bytes.NewReader(body)
	if err == nil {
		req.Body = ioutil.NopCloser(r)
	}
	return reflect.ValueOf(io.ReadSeeker(r))
}

func getCookie(req *http.Request) func(name string) *http.Cookie {
	cookies := req.Cookies()
	return func(name string) *http.Cookie {
		for _, c := range cookies {
			if c.Name == name {
				return c
			}
		}
		return nil
	}
}
