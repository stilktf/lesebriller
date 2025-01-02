package main

import (
	"log/slog"
	"net/http"
)

// This function lets you set up middlewares all routes will use.
// Taken from https://www.jvt.me/posts/2023/09/01/golang-nethttp-global-middleware/
func use(r *http.ServeMux, middlewares ...func(next http.Handler) http.Handler) http.Handler {
	var s http.Handler
	s = r

	for _, mw := range middlewares {
		s = mw(s)
	}

	return s
}

func headerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Accept") != "application/vnd.koreader.v1+json" {
			ReturnError(w, 412, 101, "Invalid Accept header format.")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	Db()
	slog.Info("initializing hashing algorithm")
	InitializeHashing()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /users/create", UserCreate)
	mux.HandleFunc("GET /users/auth", UsersAuth)

	mux.HandleFunc("PUT /syncs/progress", SyncsProgress)
	mux.HandleFunc("GET /syncs/progress/{document}", SyncsProgressPull)
	httpWrapped := use(mux, headerMiddleware)

	slog.Info("lesebriller is now running", "port", "7200")
	http.ListenAndServe(":7200", httpWrapped)
}
