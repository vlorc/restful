package advice

import (
	"github.com/vlorc/restful/pkg/binding"
	"net/http"
	"reflect"
)

func bindOf(f interface{}, c func(reflect.Type, reflect.Type) bool) {
	if v := wrap(f, c); nil != v {
		defaults = append(defaults, v)
		return
	}

	panic(errorOf(reflect.TypeOf(f)))
}

func bind(f interface{}) (reflect.Type, Advice) {
	t := reflect.TypeOf(f)
	if reflect.Func != t.Kind() || 0 == t.NumIn() {
		panic(errorOf(t))
	}

	v := reflect.ValueOf(f)
	in, on := t.NumIn(), t.NumOut()

	typ := t.In(in - 1)
	if 1 == in && on > 0 && t.Out(0) != t.In(0) && interfaceType == t.In(0) {
		return typ, bindOut(v)
	}

	call := bindCall(v)

	if 0 == on {
		return typ, func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			call(resp, req, val)
			return val
		}
	}

	if t.Out(on-1) == errorType {
		return typ, func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if r := call(resp, req, val); len(r) > 0 {
				return r[0].Interface()
			}
			return val
		}
	}

	w := New(errorType)
	return typ, func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
		if r := call(resp, req, val); len(r) > 0 {
			if r[len(r)-1].IsNil() {
				return r[0].Interface()
			}
			w(resp, req, r[len(r)-1].Interface())
		}
		return val
	}
}

func bindCall(v reflect.Value) func(http.ResponseWriter, *http.Request, interface{}) []reflect.Value {
	t := v.Type()
	b := make([]binding.Binding, t.NumIn()-1)
	for i := range b {
		b[i] = binding.New(t.In(i))
	}

	return func(resp http.ResponseWriter, req *http.Request, val interface{}) []reflect.Value {
		t := v.Type()
		p := make([]reflect.Value, len(b)+1)

		p[len(p)-1] = reflect.ValueOf(val)
		if t1, t2 := p[len(p)-1].Type(), t.In(len(p)-1); t1 != t2 {
			if !t1.ConvertibleTo(t2) {
				return nil
			}
			p[len(p)-1] = p[len(p)-1].Convert(t2)
		}
		for i := range b {
			p[i] = b[i](resp, req, t.In(i))
		}
		return v.Call(p)
	}
}

func bindOut(v reflect.Value) Advice {
	t := v.Type()
	w := make([]Advice, t.NumOut())
	for i := range w {
		w[i] = New(t.Out(i))
	}

	if t.Out(len(w)-1) != errorType {
		return func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			r := v.Call([]reflect.Value{reflect.ValueOf(val)})
			for i := range w {
				w[i](resp, req, r[i].Interface())
			}
			return val
		}
	}

	return func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
		if r := v.Call([]reflect.Value{reflect.ValueOf(val)}); r[len(r)-1].IsNil() {
			for i := range w[:len(w)-1] {
				w[i](resp, req, r[i].Interface())
			}
		} else {
			w[len(w)-1](resp, req, r[len(r)-1].Interface())
		}
		return val
	}
}

func bindFunc(t reflect.Type) Advice {
	w := make([]Advice, t.NumOut())
	for i := range w {
		w[i] = New(t.Out(i))
	}

	return func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
		if v := reflect.ValueOf(val); v.IsValid() && !v.IsZero() {
			r := v.Call(nil)
			for i := range w {
				w[i](resp, req, r[i])
			}
		}
		return val
	}
}
