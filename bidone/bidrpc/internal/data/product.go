package data

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"bidone/bidrpc/internal/biz"
)

// ProductData implements ProductRepo using in-memory storage
type ProductData struct {
	mu       sync.RWMutex
	products map[string]*biz.Product
	path     string
}

// NewProductData creates a new in-memory product repository
func NewProductData() biz.ProductRepo {
	return &ProductData{
		products: make(map[string]*biz.Product),
		path:     "./data.json",
	}
}

func (d *ProductData) get() error {
	f, err := os.ReadFile(d.path)
	if err != nil {
		return err
	}
	var records map[string]*biz.Product

	err = json.Unmarshal(f, &records)
	if err != nil {
		return err
	}
	d.products = records
	return nil
}

func (d *ProductData) set() error {

	buf, err := json.MarshalIndent(d.products, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println("--------------------------")
	if err = os.WriteFile(d.path, buf, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// Save saves a product to the in-memory store
func (d *ProductData) Save(ctx context.Context, product *biz.Product) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.products[product.ID] = product
	fmt.Println(d.set()) // save data into json file
	return nil
}

// FindByID finds a product by ID
func (d *ProductData) FindByID(ctx context.Context, id string) (*biz.Product, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.products) == 0 {
		d.get()
	}

	product, exists := d.products[id]
	if !exists {
		return nil, biz.ErrProductNotFound
	}

	// Return a copy to avoid race conditions
	return &biz.Product{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// FindAll finds all products with pagination and filtering
func (d *ProductData) FindAll(ctx context.Context, page, pageSize int32, nameFilter string) ([]*biz.Product, int32, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if len(d.products) == 0 {
		d.get()
	}

	var filtered []*biz.Product

	// Apply name filter
	for _, product := range d.products {
		if nameFilter == "" || strings.Contains(strings.ToLower(product.Name), strings.ToLower(nameFilter)) {
			filtered = append(filtered, &biz.Product{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    product.Quantity,
				CreatedAt:   product.CreatedAt,
				UpdatedAt:   product.UpdatedAt,
			})
		}
	}

	total := int32(len(filtered))

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize

	if start > total {
		return []*biz.Product{}, total, nil
	}

	if end > total {
		end = total
	}

	return filtered[start:end], total, nil
}

// Update updates an existing product
func (d *ProductData) Update(ctx context.Context, product *biz.Product) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.products) == 0 {
		d.get()
	}

	if _, exists := d.products[product.ID]; !exists {
		return biz.ErrProductNotFound
	}

	d.products[product.ID] = product
	return nil
}

// Delete deletes a product by ID
func (d *ProductData) Delete(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if len(d.products) == 0 {
		d.get()
	}

	if _, exists := d.products[id]; !exists {
		return biz.ErrProductNotFound
	}

	delete(d.products, id)
	d.set()
	return nil
}
