package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/golang/protobuf/jsonpb"

	"github.com/linzhengen/envoy-sidecar-for-grpc-lb/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(method string, err error, w http.ResponseWriter) {
	s, _ := status.FromError(err)
	if s.Code() != codes.Unimplemented {
		http.Error(w, err.Error(), 500)
		log.Printf("Something really bad happened: %v %v", s.Code(), err)
		return
	}
	http.Error(w, err.Error(), 404)
	log.Printf("Can't find %s method: %v %v", method, s.Code(), err)
}

func main() {
	grpServerHost := os.Getenv("GRPC_SERVER_HOST")
	if grpServerHost == "" {
		grpServerHost = "localhost:8080"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Printf("GRPC_SERVER_HOST is: %v", grpServerHost)
	log.Printf("PORT is: %v", port)

	conn, err := grpc.Dial(grpServerHost, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewJankenClient(conn)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	http.HandleFunc("/ping", func(w http.ResponseWriter, req *http.Request) {})

	http.HandleFunc("/janken", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Recived janken request: %v", req)
		my, _ := strconv.ParseInt(req.URL.Query().Get("my"), 10, 32)
		resp, err := c.Show(ctx, &pb.JankenRequest{My: pb.Choice(my)})
		if err != nil {
			handleError("Janken", err, w)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		m := jsonpb.Marshaler{}
		m.Marshal(w, resp)
	})
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
