package services

import (
	"context"
	"time"

	"github.com/edsjcbra/flightsimhub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductService struct {
	DB *pgxpool.Pool
}

// Construtor
func NewProductService(db *pgxpool.Pool) *ProductService {
	return &ProductService{DB: db}
}

// CreateProduct cria um novo produto
func (s *ProductService) CreateProduct(name, description string, price float64, stock int) (*models.Product, error) {
	ctx := context.Background()
	now := time.Now()
	var id int64

	err := s.DB.QueryRow(ctx,
		"INSERT INTO products (name, description, price, stock, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id",
		name, description, price, stock, now, now,
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.Product{
		ID:          id,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// GetAllProducts retorna todos os produtos
func (s *ProductService) GetAllProducts() ([]*models.Product, error) {
	ctx := context.Background()
	rows, err := s.DB.Query(ctx, "SELECT id, name, description, price, stock, created_at, updated_at FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*models.Product{}
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	return products, nil
}

// GetProductByID busca produto pelo ID
func (s *ProductService) GetProductByID(id int64) (*models.Product, error) {
	ctx := context.Background()
	var p models.Product

	err := s.DB.QueryRow(ctx,
		"SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE id=$1", id,
	).Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

// UpdateProduct atualiza um produto existente
func (s *ProductService) UpdateProduct(id int64, name, description string, price float64, stock int) (*models.Product, error) {
	ctx := context.Background()
	now := time.Now()

	_, err := s.DB.Exec(ctx,
		"UPDATE products SET name=$1, description=$2, price=$3, stock=$4, updated_at=$5 WHERE id=$6",
		name, description, price, stock, now, id,
	)
	if err != nil {
		return nil, err
	}

	return s.GetProductByID(id)
}

// DeleteProduct remove um produto pelo ID
func (s *ProductService) DeleteProduct(id int64) error {
	ctx := context.Background()
	_, err := s.DB.Exec(ctx, "DELETE FROM products WHERE id=$1", id)
	return err
}
