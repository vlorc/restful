package engine

import (
	"github.com/go-chi/chi/v5"
	"github.com/vlorc/restful/pkg/advice"
	"github.com/vlorc/restful/pkg/binding"
	"github.com/vlorc/restful/pkg/mapping"
	"github.com/vlorc/restful/pkg/render"
	"github.com/vlorc/restful/pkg/web"
	"net/http"
	"reflect"
)

func Binding() {
	binding.Register(reflect.TypeOf((*web.Param)(nil)).Elem(), func(resp http.ResponseWriter, req *http.Request, typ reflect.Type) reflect.Value {
		if ctx, ok := req.Context().Value(chi.RouteCtxKey).(chi.Context); ok {
			return reflect.ValueOf(func(k string) string {
				return ctx.URLParam(k)
			})
		}
		return reflect.ValueOf(func(k string) string {
			return ""
		})
	})
	mapping.Failed(advice.Abort, advice.Error())
}

func Advice() {
	advice.Bind(advice.Render)
	advice.Register(reflect.TypeOf((*func() interface{})(nil)).Elem(), advice.Reverse)
	advice.Register(reflect.TypeOf((*interface{})(nil)).Elem(), advice.Call)
}

func Init() {
	Binding()
	Advice()
	Rest()
}

func Rest() {
	binding.StructOf(mapping.Bind, ".+Request$", ".+Dto")
	advice.StructOf(render.Json, ".+Response$", ".+Vo", ".+Dto")
}
