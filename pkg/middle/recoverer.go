package middle

import (
	"net/http"
)

func Recover(dump func(http.ResponseWriter, *http.Request, interface{})) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if v := recover(); v != nil && v != http.ErrAbortHandler {
					dump(w, r, v)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
