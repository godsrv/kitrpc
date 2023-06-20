package http

import (
	"context"
	"encoding/json"
	"kitprc/encode"
	ep "kitprc/endpoint"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHTTPHandler(eps ep.Endpoints, opts ...kithttp.ServerOption) http.Handler {
	r := mux.NewRouter()

	r.Handle("/set", kithttp.NewServer(
		eps.PostEndpoint,
		DecodePostRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	r.Methods("GET").Path("/health").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

func DecodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ep.PostRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	return req, nil
}
