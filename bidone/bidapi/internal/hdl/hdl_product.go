package hdl

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"bidone/bidapi/internal/rpc"
	pb "bidone/bidrpc/bidrpcproto"

	chi "github.com/go-chi/chi/v5"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	var args CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		Err(w, http.StatusBadRequest, "invalid request body", err)
		return
	}
	if err := args.Validate(); err != nil {
		Err(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	req := &pb.CreateProductRequest{
		Name:        args.Name,
		Description: args.Description,
		Price:       args.Price,
		Quantity:    args.Quantity,
	}

	rsp, err := rpc.RpcClientProduct.Clt.CreateProduct(ctx, req)
	if err != nil {
		Err(w, http.StatusInternalServerError, "failed to create product", err)
		return
	}

	response := ProductDTO{
		ID:          rsp.Product.Id,
		Name:        rsp.Product.Name,
		Description: rsp.Product.Description,
		Price:       rsp.Product.Price,
		Quantity:    rsp.Product.Quantity,
		CreatedAt:   time.Unix(rsp.Product.CreatedAt, 0),
		UpdatedAt:   time.Unix(rsp.Product.UpdatedAt, 0),
	}

	Ok(w, http.StatusCreated, response)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	req := &pb.GetProductRequest{
		Id: chi.URLParam(r, "id"),
	}

	rsp, err := rpc.RpcClientProduct.Clt.GetProduct(ctx, req)
	if err != nil {
		Err(w, http.StatusNotFound, "product not found", err)
		return
	}

	response := ProductDTO{
		ID:          rsp.Product.Id,
		Name:        rsp.Product.Name,
		Description: rsp.Product.Description,
		Price:       rsp.Product.Price,
		Quantity:    rsp.Product.Quantity,
		CreatedAt:   time.Unix(rsp.Product.CreatedAt, 0),
		UpdatedAt:   time.Unix(rsp.Product.UpdatedAt, 0),
	}

	Ok(w, http.StatusOK, response)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id := chi.URLParam(r, "id")

	var args CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		Err(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	req := &pb.UpdateProductRequest{
		Id:          id,
		Name:        args.Name,
		Description: args.Description,
		Price:       args.Price,
		Quantity:    args.Quantity,
	}

	rsp, err := rpc.RpcClientProduct.Clt.UpdateProduct(ctx, req)
	if err != nil {
		Err(w, http.StatusInternalServerError, "failed to update product", err)
		return
	}

	response := ProductDTO{
		ID:          rsp.Product.Id,
		Name:        rsp.Product.Name,
		Description: rsp.Product.Description,
		Price:       rsp.Product.Price,
		Quantity:    rsp.Product.Quantity,
		CreatedAt:   time.Unix(rsp.Product.CreatedAt, 0),
		UpdatedAt:   time.Unix(rsp.Product.UpdatedAt, 0),
	}

	Ok(w, http.StatusOK, response)

}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	id := chi.URLParam(r, "id")

	req := &pb.DeleteProductRequest{
		Id: id,
	}

	rsp, err := rpc.RpcClientProduct.Clt.DeleteProduct(ctx, req)
	if err != nil {
		Err(w, http.StatusInternalServerError, "failed to delete product", err)
		return
	}
	if !rsp.Success {
		Err(w, http.StatusNotFound, "product not found", nil)
		return
	}

	Ok(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "product deleted successfully",
	})

}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 32)
	pageSize, _ := strconv.ParseInt(r.URL.Query().Get("page_size"), 10, 32)
	nameFilter := r.URL.Query().Get("name_filter")

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	req := &pb.ListProductsRequest{
		Page:       int32(page),
		PageSize:   int32(pageSize),
		NameFilter: nameFilter,
	}

	rsp, err := rpc.RpcClientProduct.Clt.ListProducts(ctx, req)
	if err != nil {
		Err(w, http.StatusInternalServerError, "failed to list products", err)
		return
	}

	productDTOs := make([]ProductDTO, len(rsp.Products))
	for i, product := range rsp.Products {
		productDTOs[i] = ProductDTO{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CreatedAt:   time.Unix(product.CreatedAt, 0),
			UpdatedAt:   time.Unix(product.UpdatedAt, 0),
		}
	}

	response := ListProductsResponse{
		Products: productDTOs,
		Total:    rsp.Total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	Ok(w, http.StatusOK, response)
}
