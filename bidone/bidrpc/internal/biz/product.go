package biz

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidInput    = errors.New("invalid input")
)

// Product represents a product in the business domain
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// ProductRepo defines the interface for product data access
type ProductRepo interface {
	Save(ctx context.Context, product *Product) error
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context, page, pageSize int32, nameFilter string) ([]*Product, int32, error)
	Update(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
}

// ProductUseCase handles product business logic
type ProductUseCase struct {
	repo ProductRepo
	mu   sync.RWMutex
}

// NewProductUseCase creates a new product use case
func NewProductUseCase(repo ProductRepo) *ProductUseCase {
	return &ProductUseCase{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (uc *ProductUseCase) CreateProduct(ctx context.Context, name, description string, price float64, quantity int32) (*Product, error) {
	if name == "" {
		return nil, ErrInvalidInput
	}
	if price < 0 {
		return nil, ErrInvalidInput
	}
	if quantity < 0 {
		return nil, ErrInvalidInput
	}

	product := &Product{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := uc.repo.Save(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

// GetProduct retrieves a product by ID
func (uc *ProductUseCase) GetProduct(ctx context.Context, id string) (*Product, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	return uc.repo.FindByID(ctx, id)
}

// ListProducts retrieves all products with pagination and filtering
func (uc *ProductUseCase) ListProducts(ctx context.Context, page, pageSize int32, nameFilter string) ([]*Product, int32, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return uc.repo.FindAll(ctx, page, pageSize, nameFilter)
}

// UpdateProduct updates an existing product
func (uc *ProductUseCase) UpdateProduct(ctx context.Context, id, name, description string, price float64, quantity int32) (*Product, error) {
	if id == "" {
		return nil, ErrInvalidInput
	}

	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		existing.Name = name
	}
	if description != "" {
		existing.Description = description
	}
	if price >= 0 {
		existing.Price = price
	}
	if quantity >= 0 {
		existing.Quantity = quantity
	}
	existing.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// DeleteProduct deletes a product by ID
func (uc *ProductUseCase) DeleteProduct(ctx context.Context, id string) error {
	if id == "" {
		return ErrInvalidInput
	}

	// Check if product exists
	_, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	return uc.repo.Delete(ctx, id)
}
