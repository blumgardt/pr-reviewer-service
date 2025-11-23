package http

import "net/http"

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}

func (r *Router) Handle(pattern string, h http.Handler) {
	r.mux.Handle(pattern, h)
}

func (r *Router) HandleFunc(pattern string, h func(w http.ResponseWriter, r *http.Request)) {
	r.mux.HandleFunc(pattern, h)
}

//func (r *Router) GET(path string, handler http.HandlerFunc) {
//
//}
