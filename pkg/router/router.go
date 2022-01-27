package router

import (
	"io/fs"
	"net/http"
)

type Router interface {
	http.Handler
	Use(middlewares ...interface{})
	Group(pattern string, fn func(Router), middlewares ...interface{})
	Get(pattern string, handles ...interface{})
	Post(pattern string, handles ...interface{})
	Put(pattern string, handles ...interface{})
	Delete(pattern string, handles ...interface{})
	Options(pattern string, handles ...interface{})
	Head(pattern string, handles ...interface{})
	Patch(pattern string, handles ...interface{})
	Trace(pattern string, handles ...interface{})
	Connect(pattern string, handles ...interface{})
	Any(pattern string, handles ...interface{})
	NotFound(handles ...interface{})
	NotMethod(handles ...interface{})
	Static(pattern string, fsys fs.FS)
	Mount(pattern string, child Router)
}
