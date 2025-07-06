package service

import (
	"context"

	pb "bidone/bidrpc/bidrpcproto"
	"bidone/bidrpc/internal/biz"
)

// ProductService implements the gRPC ProductService
type ProductService struct {
	pb.UnimplementedProductServiceServer
	uc *biz.ProductUseCase
}

// NewProductService creates a new product service
func NewProductService(uc *biz.ProductUseCase) *ProductService {
	return &ProductService{
		uc: uc,
	}
}

// CreateProduct creates a new product
func (s *ProductService) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := s.uc.CreateProduct(ctx, req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.CreateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CreatedAt:   product.CreatedAt.Unix(),
			UpdatedAt:   product.UpdatedAt.Unix(),
		},
	}, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	product, err := s.uc.GetProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CreatedAt:   product.CreatedAt.Unix(),
			UpdatedAt:   product.UpdatedAt.Unix(),
		},
	}, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	product, err := s.uc.UpdateProduct(ctx, req.Id, req.Name, req.Description, req.Price, req.Quantity)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateProductResponse{
		Product: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CreatedAt:   product.CreatedAt.Unix(),
			UpdatedAt:   product.UpdatedAt.Unix(),
		},
	}, nil
}

// DeleteProduct deletes a product by ID
func (s *ProductService) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	err := s.uc.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteProductResponse{
		Success: true,
	}, nil
}

// ListProducts lists all products with pagination and filtering
func (s *ProductService) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, total, err := s.uc.ListProducts(ctx, req.Page, req.PageSize, req.NameFilter)
	if err != nil {
		return nil, err
	}

	pbProducts := make([]*pb.Product, len(products))
	for i, product := range products {
		pbProducts[i] = &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    product.Quantity,
			CreatedAt:   product.CreatedAt.Unix(),
			UpdatedAt:   product.UpdatedAt.Unix(),
		}
	}

	return &pb.ListProductsResponse{
		Products: pbProducts,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
