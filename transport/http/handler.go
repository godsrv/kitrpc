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
		decodePostRequest,
		encode.JsonResponse,
		opts...,
	)).Methods(http.MethodPost)

	return r
}

func decodePostRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req ep.PostRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}
	return req, nil
}
