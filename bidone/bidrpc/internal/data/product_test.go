package data

import (
	"context"
	"testing"
	"time"

	"bidone/bidrpc/internal/biz"
)

func newTestProduct(id string) *biz.Product {
	return &biz.Product{
		ID:          id,
		Name:        "Test Product " + id,
		Description: "A test product",
		Price:       10.0,
		Quantity:    5,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// Test FindAll pagination and filter
func TestProductData_FindAll_PaginationAndFilter(t *testing.T) {
	d := NewProductData()
	ctx := context.Background()
	// Add multiple products
	for i := 1; i <= 15; i++ {
		p := newTestProduct("p" + string(rune(i+'0')))
		p.Name = "Product" + string(rune(i+'0'))
		if i%2 == 0 {
			p.Name = "Special" + p.Name
		}
		if err := d.Save(ctx, p); err != nil {
			t.Fatalf("Save failed: %v", err)
		}
	}

	// Test filter for "Special" and pagination (page 1, pageSize 3)
	products, total, err := d.FindAll(ctx, 1, 3, "Special")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 7 {
		t.Errorf("Expected 7 products with 'Special' in name, got %d", total)
	}
	if len(products) != 3 {
		t.Errorf("Expected 3 products on page 1, got %d", len(products))
	}

	// Test page 3 (should have 1 product)
	products, _, err = d.FindAll(ctx, 3, 3, "Special")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if len(products) != 1 {
		t.Errorf("Expected 1 product on page 3, got %d", len(products))
	}
}

func TestProductData_CRUD(t *testing.T) {
	d := NewProductData()
	ctx := context.Background()
	p := newTestProduct("p1")

	// Test Save
	if err := d.Save(ctx, p); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Test FindByID
	got, err := d.FindByID(ctx, p.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("FindByID returned wrong product: got %v, want %v", got.ID, p.ID)
	}

	// Test Update
	p.Name = "Updated Name"
	if err := d.Update(ctx, p); err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	got, _ = d.FindByID(ctx, p.ID)
	if got.Name != "Updated Name" {
		t.Errorf("Update did not update name: got %v", got.Name)
	}

	// Test FindAll
	ps, total, err := d.FindAll(ctx, 1, 10, "Updated")
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}
	if total != 1 || len(ps) != 1 {
		t.Errorf("FindAll returned wrong count: got %d, want 1", total)
	}

	// Test Delete
	if err := d.Delete(ctx, p.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	_, err = d.FindByID(ctx, p.ID)
	if err == nil {
		t.Errorf("FindByID should fail after delete")
	}
}
