package main

import (
	"context"
	"fmt"
	"kitprc/transport/grpc/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:7002", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	svc := pb.NewServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := svc.Post(ctx, &pb.PostRequest{
		Key: "hello",
		Val: "world",
	})
	if err != nil {
		log.Fatalf("could not post: %v", err)
	}

	fmt.Println("data: %s", r.GetData())
}
