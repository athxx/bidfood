package main

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/athxx/bidfood/bidapi/internal/hdl"
	"github.com/athxx/bidfood/bidapi/internal/rpc"

	"github.com/go-chi/chi/v5"
)

var r = NewTestRouter()
var prodID = `82deb197-34e4-4e28-add3-560a81da47cb`

// mock  router returns a http.Handler for testing
func NewTestRouter() http.Handler {
	grpcAddr := "localhost:9000" // grpc server address
	if err := rpc.InitProductGrpcClient(grpcAddr); err != nil {
		log.Fatalf("failed to initialize product gRPC client: %v", err)
	}
	r := chi.NewRouter()
	r.Post("/products", hdl.CreateProduct)
	r.Get("/products", hdl.ListProducts)
	r.Get("/products/{id}", hdl.GetProduct)
	r.Put("/products/{id}", hdl.UpdateProduct)
	r.Delete("/products/{id}", hdl.DeleteProduct)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	return r
}

func TestHealthCheck(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/health")
	if err != nil {
		t.Fatalf("GET /health failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if string(body) != "OK" {
		t.Errorf("Expected body 'OK', got '%s'", string(body))
	}
}

func TestCreateProduct_Success(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/products", "application/json", strings.NewReader(`{
		"name": "iPhone 15 Pro",
		"description": "Latest iPhone with advanced camera system",
		"price": 999.99,
		"quantity": 50}`))
	if err != nil {
		t.Fatalf("POST /products failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected error status, got 200")
	}
}

func TestCreateProduct_BadRequest(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/products", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatalf("POST /products failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		t.Errorf("Expected error status, got 200")
	}
}

func TestListProduct(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/products")
	if err != nil {
		t.Fatalf("POST /products failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 error status, got 200")
	}
}

func TestListProduct_Paginator(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/products?page=1&page_size=1")
	if err != nil {
		t.Fatalf("POST /products failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 error status, got 200")
	}
}

func TestGetProduct(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/products/" + prodID)
	if err != nil {
		t.Fatalf("GET /products/%s failed: %v", prodID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200=, got %d", resp.StatusCode)
	}
}

func TestUpdateProduct(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/products/"+prodID, strings.NewReader(`{"name":"Updated"}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("PUT /products/%s failed: %v", prodID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for update, got %d", resp.StatusCode)
	}
}

func TestDeleteProduct(t *testing.T) {
	ts := httptest.NewServer(r)
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/products/"+prodID, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("DELETE /products/%s failed: %v", prodID, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 for delete, got %d", resp.StatusCode)
	}
}
