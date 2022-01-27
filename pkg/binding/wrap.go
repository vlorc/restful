package binding

import (
	"net/http"
	"reflect"
)

func Wrap(f interface{}, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Binding {
	return wrap(f, c)
}

func wrap(f interface{}, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Binding {
	switch v := f.(type) {
	case func(reflect.Type) Binding:
		return v
	case func(http.ResponseWriter, *http.Request, reflect.Type) reflect.Value:
		return wrapInterface(v, c)
	case Binding:
		return wrapInterface(v, c)
	}

	if t, v := bind(f); nil != t && nil != v {
		return wrapType(v, t, c)
	}

	return nil
}

func wrapType(v Binding, t reflect.Type, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Binding {
	return func(r reflect.Type) Binding {
		if !c(t, r) {
			return nil
		}
		if t == r {
			return v
		}
		return func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
			if rv := v(resp, req, typ); rv.Type() == r {
				return rv
			} else if rv.Type().ConvertibleTo(r) {
				return rv.Convert(r)
			}
			return reflect.Zero(r)
		}
	}
}

func wrapInterface(v Binding, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Binding {
	return func(r reflect.Type) Binding {
		if c(interfaceType, r) {
			return v
		}
		return nil
	}
}
