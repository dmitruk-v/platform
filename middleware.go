package platform

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

func CorsMiddleware(log *log.Logger, allowed ...string) middleware {
	return func(h handler) handler {
		wrapper := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var isAllowed bool
			origin := r.Header.Get("Origin")
			for _, client := range allowed {
				if origin == client {
					isAllowed = true
					break
				}
			}
			if isAllowed {
				w.Header().Add("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
			if r.Method == http.MethodOptions {
				w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, HEAD")
				w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
				// w.Header().Add("Access-Control-Max-Age", "86400")
			}
			if err := h(ctx, w, r); err != nil {
				return err
			}
			return nil
		}
		return wrapper
	}
}

func ErrorsMiddleware(log *log.Logger) middleware {
	return func(h handler) handler {
		wrapper := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := h(ctx, w, r); err != nil {
				log.Printf("ERROR %v -> %v\n", r.RemoteAddr, err)
			}
			if ctx.Err() != nil {
				log.Printf("ERROR %v -> %v\n", r.RemoteAddr, ctx.Err())
			}
			return nil
		}
		return wrapper
	}
}

func TraceMiddleware(log *log.Logger) middleware {
	return func(h handler) handler {
		wrapper := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			start := time.Now()
			if err := h(ctx, w, r); err != nil {
				return fmt.Errorf("trace middleware, with error: %v", err)
			}

			val, ok := ctx.Value(AppKey).(*CtxValue)
			if !ok {
				return ErrContextValue
			}

			spent := time.Since(start)
			log.Printf("TRACE %v ( %v ) -> %v %v %v, %v", r.RemoteAddr, r.Header.Get("Origin"), r.Method, r.RequestURI, val.StatusCode, spent)
			return nil
		}
		return wrapper
	}
}
