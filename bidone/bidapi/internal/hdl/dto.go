package hdl

import (
	"errors"
	"time"
)

type ProductDTO struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int32   `json:"quantity"`
}

type ListProductsResponse struct {
	Products []ProductDTO `json:"products"`
	Total    int32        `json:"total"`
	Page     int32        `json:"page"`
	PageSize int32        `json:"page_size"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (p *CreateProductRequest) Validate() error {
	if p.Name == "" {
		return errors.New("Name is required")
	}
	if p.Price <= 0 {
		return errors.New("Price must be greater than zero")
	}
	if p.Quantity < 0 {
		return errors.New("Quantity must be greater than zero")
	}
	return nil
}
