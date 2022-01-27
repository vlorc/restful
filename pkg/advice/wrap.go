package advice

import (
	"net/http"
	"reflect"
)

func Wrap(f interface{}, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Advice {
	return wrap(f, c)
}

func wrap(f interface{}, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Advice {
	switch v := f.(type) {
	case func() Advice:
		return func(reflect.Type) Advice {
			return v()
		}
	case func(reflect.Type) Advice:
		return v
	case func(http.ResponseWriter, *http.Request, interface{}) interface{}:
		return wrapInterface(v, c)
	case Advice:
		return wrapInterface(v, c)
	case func(interface{}) interface{}:
		return wrapInterface(func(resp http.ResponseWriter, req *http.Request, val interface{}) interface{} {
			if nil != val {
				val = v(val)
			}
			return val
		}, c)
	}

	if t, v := bind(f); nil != t && nil != v {
		return wrapType(v, t, c)
	}

	return nil
}

func wrapInterface(v Advice, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Advice {
	return func(r reflect.Type) Advice {
		if c(interfaceType, r) {
			return v
		}
		return nil
	}
}

func wrapType(v Advice, t reflect.Type, c func(reflect.Type, reflect.Type) bool) func(reflect.Type) Advice {
	return func(r reflect.Type) Advice {
		if c(t, r) {
			return v
		}
		return nil
	}
}
