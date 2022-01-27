package mapping

import (
	"github.com/vlorc/restful/pkg/web"
	"net/http"
	"reflect"
)

const (
	MIME_JSON              = "application/json"
	MIME_HTML              = "text/html"
	MIME_XML               = "application/xml"
	MIME_XML2              = "text/xml"
	MIME_Plain             = "text/plain"
	MIME_POSTForm          = "application/x-www-form-urlencoded"
	MIME_MultipartPOSTForm = "multipart/form-data"
	MIME_PROTOBUF          = "application/x-protobuf"
	MIME_MSGPACK           = "application/x-msgpack"
	MIME_MSGPACK2          = "application/msgpack"
	MIME_YAML              = "application/x-yaml"
)

var (
	JSON          = jsonMapping{}
	XML           = xmlMapping{}
	Form          = formMapping{}
	Query         = queryMapping{}
	FormPost      = formPostMapping{}
	FormMultipart = formMultipartMapping{}
	Header        = headerMapping{}
)

type Mapping interface {
	Name() string
	Bind(*http.Request, interface{}) error
}

var defaults []func(*http.Request, reflect.Type) Mapping
var success []func(http.ResponseWriter, *http.Request, interface{}) interface{}
var failed []func(http.ResponseWriter, *http.Request, interface{}) interface{}

var mapping = map[string]Mapping{
	MIME_JSON:              JSON,
	MIME_XML:               XML,
	MIME_XML2:              XML,
	MIME_POSTForm:          FormPost,
	MIME_MultipartPOSTForm: FormMultipart,
}

func Register(t string, m Mapping) {
	if nil != mapping {
		mapping[t] = m
	} else {
		delete(mapping, t)
	}
}

func Default(m func(*http.Request, reflect.Type) Mapping) {
	defaults = append(defaults, m)
}

func Success(f ...func(http.ResponseWriter, *http.Request, interface{}) interface{}) {
	if len(f) > 0 {
		success = append(success, f...)
	}
}

func Failed(f ...func(http.ResponseWriter, *http.Request, interface{}) interface{}) {
	if len(f) > 0 {
		failed = append(failed, f...)
	}
}

func Load(req *http.Request, typ reflect.Type) Mapping {
	t := filterFlags(req.Header.Get(web.HTTP_HEADER_CONTENT_TYPE))
	if m := mapping[t]; nil != m {
		return m
	}

	for i := len(defaults) - 1; i >= 0; i-- {
		if m := defaults[i](req, typ); nil != m {
			return m
		}
	}

	return Form
}

func Bind(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
	m := Load(req, typ)

	var v reflect.Value
	var r interface{}

	if reflect.Ptr == typ.Kind() {
		v = reflect.New(typ.Elem())
		r = v.Interface()
	} else {
		v = reflect.New(typ)
		r = v.Interface()
		v = v.Elem()
	}

	err := m.Bind(req, r)
	if nil == err {
		if err = doSuccess(resp, req, r); nil == err {
			return v
		}
	}
	doFailed(resp, req, err)

	panic(http.ErrAbortHandler)
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

func doSuccess(resp http.ResponseWriter, req *http.Request, val interface{}) error {
	for i := len(success) - 1; i >= 0; i-- {
		if err, ok := success[i](resp, req, val).(error); ok && nil != err {
			return err
		}
	}
	return nil
}

func doFailed(resp http.ResponseWriter, req *http.Request, err error) {
	var r interface{} = err
	for i := len(failed) - 1; i >= 0; i-- {
		if r = failed[i](resp, req, r); nil == r {
			break
		}
	}
}
