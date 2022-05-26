package chi

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vlorc/restful/pkg/engine"
	"github.com/vlorc/restful/pkg/engine/core"
	"github.com/vlorc/restful/pkg/middle"
	"github.com/vlorc/restful/pkg/router"
	"github.com/vlorc/restful/pkg/web"
	"github.com/vlorc/restful/pkg/wrap"
	"io/fs"
	"net/http"
	"sync/atomic"
)

const CHI_NAME = "chi"

type Router struct {
	router   chi.Router
	pattern  func(string) string
	middles  []func(http.Handler) http.Handler
	state    int32
	middleOf func(interface{}) func(http.Handler) http.Handler
	wrapOf   func(...interface{}) http.HandlerFunc
}

var _ router.Router = &Router{}

func init() {
	engine.Register(engine.DETAULT_NAME, func(middles ...func(http.Handler) http.Handler) router.Router {
		core.Setup(chi.RouteCtxKey)
		return NewRouter(middles...)
	})
	engine.Register(CHI_NAME, func(middles ...func(http.Handler) http.Handler) router.Router {
		core.Setup(chi.RouteCtxKey)
		return NewRouter(middles...)
	})
}

func NewRouter(middles ...func(http.Handler) http.Handler) router.Router {
	r := &Router{
		router:   chi.NewRouter(),
		pattern:  router.Pattern(""),
		middleOf: middle.Middle,
		wrapOf:   wrap.Wrap,
		middles:  middles,
	}

	r.middles = append(r.middles, middle.Recover(func(_ http.ResponseWriter, _ *http.Request, val interface{}) {
		middleware.PrintPrettyStack(val)
	}), r.middleOf(web.NewResponseWriter))

	return r
}

func (r *Router) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(resp, req)
}

func (r *Router) Use(middlewares ...interface{}) {
	if atomic.LoadInt32(&r.state) == 0 {
		for i := range middlewares {
			m := r.middleOf(middlewares[i])
			r.middles = append(r.middles, m)
		}
	}
}

func (r *Router) Group(pattern string, fn func(router.Router), middlewares ...interface{}) {
	r.setup()

	f := func(rc chi.Router) {
		rr := &Router{
			router:   rc,
			pattern:  router.Pattern(""),
			wrapOf:   r.wrapOf,
			middleOf: r.middleOf,
		}
		rr.Use(middlewares)
		fn(rr)
	}

	if "" != pattern {
		p := r.pattern(pattern)
		r.router.Route(p, f)
	} else {
		r.router.Group(f)
	}
}

func (r *Router) Get(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Get(p, r.wrapOf(handles...))
	}
}

func (r *Router) Post(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Post(p, r.wrapOf(handles...))
	}
}

func (r *Router) Put(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Put(p, r.wrapOf(handles...))
	}
}

func (r *Router) Delete(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Delete(p, r.wrapOf(handles...))
	}
}

func (r *Router) Options(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Options(p, r.wrapOf(handles...))
	}
}

func (r *Router) Head(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Head(p, r.wrapOf(handles...))
	}
}

func (r *Router) Patch(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Patch(p, r.wrapOf(handles...))
	}
}

func (r *Router) Connect(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Connect(p, r.wrapOf(handles...))
	}
}

func (r *Router) Trace(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.Trace(p, r.wrapOf(handles...))
	}
}

func (r *Router) Any(pattern string, handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		p := r.pattern(pattern)
		r.router.HandleFunc(p, r.wrapOf(handles...))
	}
}

func (r *Router) NotFound(handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		r.router.NotFound(r.wrapOf(handles...))
	}
}

func (r *Router) NotMethod(handles ...interface{}) {
	r.setup()

	if len(handles) > 0 {
		r.router.MethodNotAllowed(r.wrapOf(handles...))
	}
}

func (r *Router) Static(pattern string, fsys fs.FS) {
	r.setup()

	p := r.pattern(pattern)
	r.router.Handle(p, http.FileServer(http.FS(fsys)))
}

func (r *Router) Mount(pattern string, child router.Router) {
	r.setup()

	r.router.Mount(pattern, child)
}

func (r *Router) setup() {
	if atomic.CompareAndSwapInt32(&r.state, 0, 1) && len(r.middles) > 0 {
		v := make([]func(http.Handler) http.Handler, len(r.middles))
		copy(v, r.middles)
		for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
			v[i], v[j] = v[j], v[i]
		}
		r.router.Use(v...)
	}
}
