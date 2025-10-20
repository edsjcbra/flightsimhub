package services

import (
	"context"
	"time"

	"github.com/edsjcbra/flightsimhub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartService struct {
	DB *pgxpool.Pool
}

// Novo construtor usando *pgxpool.Pool diretamente
func NewCartService(db *pgxpool.Pool) *CartService {
	return &CartService{DB: db}
}

// AddItem adiciona um produto ao carrinho do usuário
func (s *CartService) AddItem(userID, productID, quantity int) error {
	ctx := context.Background()

	// Verifica se o usuário já tem um carrinho
	var cartID int
	err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID)
	if err != nil {
		// Cria carrinho se não existir
		err = s.DB.QueryRow(ctx, "INSERT INTO carts (user_id, created_at, updated_at) VALUES ($1, $2, $3) RETURNING id",
			userID, time.Now(), time.Now()).Scan(&cartID)
		if err != nil {
			return err
		}
	}

	// Insere item no carrinho
	_, err = s.DB.Exec(ctx,
		`INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		cartID, productID, quantity, time.Now(), time.Now())
	return err
}

// GetCart retorna todos os itens do carrinho do usuário
func (s *CartService) GetCart(userID int) ([]models.CartItem, error) {
	ctx := context.Background()

	// Pega o carrinho do usuário
	var cartID int
	err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID)
	if err != nil {
		return nil, err
	}

	rows, err := s.DB.Query(ctx, "SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id=$1", cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		err := rows.Scan(&item.ID, &item.CartID, &item.ProductID, &item.Quantity, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

// RemoveItem remove um item do carrinho
func (s *CartService) RemoveItem(itemID int) error {
	ctx := context.Background()
	_, err := s.DB.Exec(ctx, "DELETE FROM cart_items WHERE id=$1", itemID)
	return err
}
