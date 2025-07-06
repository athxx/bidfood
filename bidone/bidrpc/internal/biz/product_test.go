package biz

import (
	"context"
	"testing"
)

type mockProductRepo struct {
	products map[string]*Product
}

func newMockProductRepo() *mockProductRepo {
	return &mockProductRepo{products: make(map[string]*Product)}
}

func (m *mockProductRepo) Save(ctx context.Context, product *Product) error {
	m.products[product.ID] = product
	return nil
}
func (m *mockProductRepo) FindByID(ctx context.Context, id string) (*Product, error) {
	p, ok := m.products[id]
	if !ok {
		return nil, ErrProductNotFound
	}
	return p, nil
}
func (m *mockProductRepo) FindAll(ctx context.Context, page, pageSize int32, nameFilter string) ([]*Product, int32, error) {
	var out []*Product
	for _, p := range m.products {
		out = append(out, p)
	}
	total := int32(len(out))
	return out, total, nil
}
func (m *mockProductRepo) Update(ctx context.Context, product *Product) error {
	if _, ok := m.products[product.ID]; !ok {
		return ErrProductNotFound
	}
	m.products[product.ID] = product
	return nil
}
func (m *mockProductRepo) Delete(ctx context.Context, id string) error {
	if _, ok := m.products[id]; !ok {
		return ErrProductNotFound
	}
	delete(m.products, id)
	return nil
}

func TestProductUseCase_CRUD(t *testing.T) {
	repo := newMockProductRepo()
	uc := NewProductUseCase(repo)
	ctx := context.Background()

	// Create
	p, err := uc.CreateProduct(ctx, "name", "desc", 1.2, 3)
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}
	if p.ID == "" {
		t.Error("CreateProduct did not set ID")
	}

	// Get
	got, err := uc.GetProduct(ctx, p.ID)
	if err != nil {
		t.Fatalf("GetProduct failed: %v", err)
	}
	if got.Name != "name" {
		t.Errorf("GetProduct wrong name: got %v", got.Name)
	}

	// List
	list, total, err := uc.ListProducts(ctx, 1, 10, "")
	if err != nil || total != 1 || len(list) != 1 {
		t.Errorf("ListProducts failed: err=%v, total=%d, len=%d", err, total, len(list))
	}

	// Update
	updated, err := uc.UpdateProduct(ctx, p.ID, "newname", "newdesc", 2.3, 5)
	if err != nil {
		t.Fatalf("UpdateProduct failed: %v", err)
	}
	if updated.Name != "newname" || updated.Price != 2.3 {
		t.Errorf("UpdateProduct did not update fields")
	}

	// Delete
	if err := uc.DeleteProduct(ctx, p.ID); err != nil {
		t.Fatalf("DeleteProduct failed: %v", err)
	}
	_, err = uc.GetProduct(ctx, p.ID)
	if err == nil {
		t.Error("GetProduct should fail after delete")
	}
}
