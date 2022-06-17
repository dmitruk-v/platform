package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PayloadResponse struct {
	Code    int         `json:"code"`
	Payload interface{} `json:"payload"`
}

func RespondError(ctx context.Context, w http.ResponseWriter, res ErrorResponse) error {
	return respond(ctx, w, res)
}

func RespondPayload(ctx context.Context, w http.ResponseWriter, res PayloadResponse) error {
	return respond(ctx, w, res)
}

func respond(ctx context.Context, w http.ResponseWriter, res interface{}) error {
	ctxVal, ok := ctx.Value(AppKey).(*CtxValue)
	if !ok {
		return ErrContextValue
	}
	w.Header().Add("Content-Type", "application/json")
	switch r := res.(type) {
	case ErrorResponse:
		ctxVal.StatusCode = r.Code
		w.WriteHeader(r.Code)
	case PayloadResponse:
		ctxVal.StatusCode = r.Code
		w.WriteHeader(r.Code)
	default:
		return ErrResponseType
	}
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return fmt.Errorf("encoding response payload to json, with error: %v", err)
	}
	return nil
}
