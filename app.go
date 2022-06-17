package platform

import (
	"context"
	"log"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

type handler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error
type middleware func(h handler) handler

type CtxKey string
type CtxValue struct {
	StatusCode int
	Message    string
}

const AppKey CtxKey = "qwerty12345"

type app struct {
	mux *httptreemux.ContextMux
	mws []middleware
}

func NewApp() *app {
	return &app{
		mux: httptreemux.NewContextMux(),
		mws: make([]middleware, 0),
	}
}

func (a *app) wrapMiddleware(h handler, mws []middleware) http.HandlerFunc {
	var wrapper handler = h
	for i := len(mws) - 1; i >= 0; i-- {
		mw := mws[i]
		wrapper = mw(wrapper)
	}
	hand := func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), AppKey, &CtxValue{})
		if err := wrapper(ctx, w, r); err != nil {
			log.Println(err)
		}
	}
	return hand
}

func (a *app) Middleware(mw ...middleware) {
	a.mws = append(a.mws, mw...)
}

func (a *app) Handle(method string, pattern string, h handler, mws ...middleware) {
	if len(mws) > 0 {
		// Use handler-specific middleware
		a.mux.Handle(method, pattern, a.wrapMiddleware(h, mws))
	} else {
		// Use global app middleware
		a.mux.Handle(method, pattern, a.wrapMiddleware(h, a.mws))
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.mux.ServeHTTP(w, r)
}
