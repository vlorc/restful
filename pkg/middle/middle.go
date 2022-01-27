package middle

import "net/http"

func Middle(middlewares interface{}) func(http.Handler) http.Handler {
	switch v := middlewares.(type) {
	case func(http.Handler) http.Handler:
		return v
	case func(http.ResponseWriter, *http.Request):
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				v(resp, req)
				next.ServeHTTP(resp, req)
			})
		}
	case func(*http.Request) *http.Request:
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				next.ServeHTTP(resp, v(req))
			})
		}
	case func(http.ResponseWriter) http.ResponseWriter:
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				next.ServeHTTP(v(resp), req)
			})
		}
	case http.Handler:
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
				v.ServeHTTP(resp, req)
				next.ServeHTTP(resp, req)
			})
		}
	}

	return func(next http.Handler) http.Handler {
		return next
	}
}
