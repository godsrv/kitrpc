package grpc

import (
	"context"
	"encoding/json"

	"kitprc/encode"
	ep "kitprc/endpoint"
	"kitprc/transport/grpc/pb"

	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	post grpc.Handler
	pb.UnimplementedServiceServer
}

func (g *grpcServer) Post(ctx context.Context, req *pb.PostRequest) (*pb.ServiceResponse, error) {
	_, rep, err := g.post.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceResponse), nil
}

func MakeGRPCHandler(eps ep.Endpoints, opts ...grpc.ServerOption) pb.ServiceServer {
	return &grpcServer{
		post: grpc.NewServer(
			eps.PostEndpoint,
			decodePostRequest,
			encodeResponse,
			opts...,
		),
	}
}

func decodePostRequest(_ context.Context, r interface{}) (interface{}, error) {
	return ep.PostRequest{
		Key: r.(*pb.PostRequest).Key,
		Val: r.(*pb.PostRequest).Val,
	}, nil
}

func encodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(encode.Response)
	data, _ := json.Marshal(resp.Data)
	return &pb.ServiceResponse{
		Code:    int64(resp.Code),
		Data:    string(data),
		Message: "success",
	}, nil
}
