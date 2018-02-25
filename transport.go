package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

type addRequest struct {
	V int `json:"v"`
}

func (r addRequest) String() string {
	return fmt.Sprintf("%d", r.V)
}

type addResponse struct {
	V int `json:"v"`
}

func (r addResponse) String() string {
	return fmt.Sprintf("%d", r.V)
}

func makeAddEndpoint(svc Counter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		v := svc.Add(req.V)
		return addResponse{v}, nil
	}
}

func decodeAddRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req addRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
