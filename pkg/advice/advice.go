package advice

import (
	"fmt"
	"net/http"
	"github.com/vlorc/restful/pkg/web"
	"reflect"
)

type Advice func(http.ResponseWriter, *http.Request, interface{}) interface{}

var mapping = types()
var defaults []func(reflect.Type) Advice
var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func Register(t reflect.Type, o Advice) {
	if nil != t {
		if nil != o {
			mapping[t] = append(mapping[t], o)
		} else {
			delete(mapping, t)
		}
	}
}

func Bind(f interface{}) {
	bindOf(f, func(t reflect.Type, r reflect.Type) bool {
		return r.ConvertibleTo(t)
	})
}

func BindOf(f interface{}, c func(reflect.Type, reflect.Type) bool) {
	bindOf(f, c)
}

func Empty() Advice {
	return empty
}

func StructOf(f interface{}, s ...string) {
	m := web.Match(s...)
	bindOf(f, func(_ reflect.Type, r reflect.Type) bool {
		r = elem(r)
		return reflect.Struct == r.Kind() && m(r.String())
	})
}

func Call(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
	if nil != val {
		if err, ok := New(reflect.TypeOf(val))(resp, req, val).(error); !ok || nil != err {
			return nil
		}
	}
	return val
}

func New(t reflect.Type) Advice {
	if o, ok := mapping[t]; ok {
		return Group(o...)
	}

	if reflect.Func == t.Kind() && 0 == t.NumIn() && t.NumOut() > 0 {
		if o := bindFunc(t); nil != o {
			return o
		}
	}

	for i := len(defaults) - 1; i >= 0; i-- {
		if o := defaults[i](t); nil != o {
			return o
		}
	}

	panic(errorOf(t))
}

func Group(a ...Advice) Advice {
	if len(a) == 0 {
		return empty
	}
	if len(a) == 1 {
		return a[0]
	}

	return func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
		v := val
		for i := len(a) - 1; i >= 0; i-- {
			if v = a[i](resp, req, v); nil == v {
				break
			}
		}
		return v
	}
}

func errorOf(t reflect.Type) error {
	return fmt.Errorf("Unsupported advice type: %s", t.String())
}

func elem(r reflect.Type) reflect.Type {
	if reflect.Ptr == r.Kind() {
		return r.Elem()
	}
	return r
}

func empty(_ http.ResponseWriter, _ *http.Request, v interface{}) interface{} {
	return v
}
