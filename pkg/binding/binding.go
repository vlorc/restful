package binding

import (
	"fmt"
	"net/http"
	"github.com/vlorc/restful/pkg/web"
	"reflect"
)

type Binding func(http.ResponseWriter, *http.Request, reflect.Type) reflect.Value

var mapping = types()
var defaults []func(reflect.Type) Binding
var interfaceType = reflect.TypeOf((*interface{})(nil)).Elem()

func Register(t reflect.Type, b Binding) {
	if nil != t {
		if nil != b {
			mapping[t] = b
		} else {
			delete(mapping, t)
		}
	}
}

func Zero(reflect.Type) Binding {
	return func(_ http.ResponseWriter, _ *http.Request, typ reflect.Type) reflect.Value {
		return reflect.Zero(typ)
	}
}

func New(t reflect.Type) Binding {
	if b, ok := mapping[t]; ok {
		return b
	}

	for i := len(defaults) - 1; i >= 0; i-- {
		if b := defaults[i](t); nil != b {
			return b
		}
	}

	panic(errorOf(t))
}

func Bind(f interface{}) {
	bindOf(f, func(t reflect.Type, r reflect.Type) bool {
		return r.ConvertibleTo(t)
	})
}

func BindOf(f interface{}, c func(reflect.Type, reflect.Type) bool) {
	bindOf(f, c)
}

func errorOf(t reflect.Type) error {
	return fmt.Errorf("Unsupported binding type: %s", t.String())
}

func bindOf(f interface{}, c func(reflect.Type, reflect.Type) bool) {
	if v := wrap(f, c); nil != v {
		defaults = append(defaults, v)
		return
	}

	panic(errorOf(reflect.TypeOf(f)))
}

func StructOf(f interface{}, s ...string) {
	m := web.Match(s...)
	bindOf(f, func(_ reflect.Type, r reflect.Type) bool {
		r = elem(r)
		return reflect.Struct == r.Kind() && m(r.String())
	})
}

func isAnonymous(t reflect.Type) bool {
	t = elem(t)
	return reflect.Struct == t.Kind() && "" == t.Name()
}

func bind(f interface{}) (reflect.Type, Binding) {
	t := reflect.TypeOf(f)
	if reflect.Func != t.Kind() || 0 == t.NumOut() {
		panic(errorOf(t))
	}

	b := make([]Binding, t.NumIn())
	for i := range b {
		b[i] = New(t.In(i))
	}

	v := reflect.ValueOf(f)

	return t.Out(0), func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
		p := make([]reflect.Value, len(b))
		for i := range p {
			p[i] = b[i](resp, req, typ)
		}
		return v.Call(p)[0]
	}
}

func elem(r reflect.Type) reflect.Type {
	if reflect.Ptr == r.Kind() {
		return r.Elem()
	}
	return r
}
