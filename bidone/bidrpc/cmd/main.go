package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "bidone/bidrpc/bidrpcproto"
	"bidone/bidrpc/internal/biz"
	"bidone/bidrpc/internal/data"
	"bidone/bidrpc/internal/service"

	"google.golang.org/grpc"
)

var (
	port = flag.String("port", "9000", "gRPC server port")
)

func main() {
	flag.Parse()

	// Initialize repository
	repo := data.NewProductData()

	// Initialize use case
	uc := biz.NewProductUseCase(repo)

	// Initialize service
	productService := service.NewProductService(uc)

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterProductServiceServer(s, productService)

	// Start server
	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("gRPC server listening on port %s", *port)

	// graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c

		log.Println("Shutting down gRPC server...")
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
