package main

import (
	"context"
	"log"
	"math/rand"
	"net"
	"os"

	"github.com/linzhengen/envoy-sidecar-for-grpc-lb/pb"
	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct{}

// Show the hand
func (s *server) Show(ctx context.Context, in *pb.JankenRequest) (*pb.JankenResponse, error) {
	log.Printf("Handling JankenRequest request [%v] with context %v", in, ctx)
	hostname, err := os.Hostname()
	if err != nil {
		log.Printf("Unable to get hostname %v", err)
		hostname = ""
	}
	grpc.SendHeader(ctx, metadata.Pairs("hostname", hostname))
	c := choice(pb.Choice_name)
	log.Println(c)
	return &pb.JankenResponse{
		My:    in.My,
		Your:  c,
		IsWin: isWin(in.My, c),
	}, nil
}

func choice(m map[int32]string) pb.Choice {
	return pb.Choice(int32(rand.Intn(len(m))))
}

func isWin(o, t pb.Choice) bool {
	return (o == pb.Choice_ROCK && t == pb.Choice_PAPER) || (o == pb.Choice_PAPER && t == pb.Choice_SCISSORS) || (o == pb.Choice_SCISSORS && t == pb.Choice_ROCK)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterJankenServer(grpcServer, &server{})
	grpc_health_v1.RegisterHealthServer(grpcServer, &health.Server{})
	reflection.Register(grpcServer)
	log.Printf("Listening for Janken on port %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
