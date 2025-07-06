package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/athxx/bidfood/bidapi/internal/hdl"
	"github.com/athxx/bidfood/bidapi/internal/rpc"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	port     = flag.String("port", "8080", "HTTP server port")
	grpcAddr = flag.String("grpc-addr", "localhost:9000", "gRPC server address")
)

func init() {
	flag.Parse()
	if err := rpc.InitProductGrpcClient(*grpcAddr); err != nil {
		log.Fatalf("failed to initialize product gRPC client: %v", err)
	}
}

func main() {
	// Close product rpc client
	defer rpc.RpcClientProduct.Close()

	// Create HTTP router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// r.Route("/api/v1", func(r chi.Router) {
	// 	// Middleware for API versioning

	// Product routes
	r.Post("/products", hdl.CreateProduct)
	r.Get("/products", hdl.ListProducts)
	r.Get("/products/{id}", hdl.GetProduct)
	r.Put("/products/{id}", hdl.UpdateProduct)
	r.Delete("/products/{id}", hdl.DeleteProduct)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + *port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		<-c
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Printf("HTTP server listening on port %s", *port)

	// Start server
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	<-ctx.Done()

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("Shutting down HTTP server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP server shutdown error: %v", err)
	}
}
