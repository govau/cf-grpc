package main

//go:generate protoc -I ../ ../service.proto --go_out=plugins=grpc:../pb

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/govau/cf-grpc/pb"
)

type backendServer struct {
	Messages map[string]string
}

func (s *backendServer) GetStatus(ctx context.Context, req *pb.StatusRequest) (*pb.StatusResponse, error) {
	resp, ok := s.Messages[req.Job]
	if !ok {
		return nil, errors.New("no such job")
	}
	return &pb.StatusResponse{
		Status: resp,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBackendServer(grpcServer, &backendServer{
		Messages: map[string]string{
			"foo":  "bar",
			"biz:": "boz",
		},
	})
	log.Fatal(grpcServer.Serve(lis))
}
