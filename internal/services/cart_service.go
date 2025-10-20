package services

import (
	"context"
	"errors"
	"time"

	"github.com/edsjcbra/flightsimhub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartService struct {
	DB *pgxpool.Pool
}

func NewCartService(db *pgxpool.Pool) *CartService {
	return &CartService{DB: db}
}

// ensureCart returns cartID for user (creates if not exists)
func (s *CartService) ensureCart(ctx context.Context, userID int) (int, error) {
	var cartID int
	err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID)
	if err == nil {
		return cartID, nil
	}
	// create
	err = s.DB.QueryRow(ctx, "INSERT INTO carts (user_id, created_at, updated_at) VALUES ($1,$2,$3) RETURNING id",
		userID, time.Now(), time.Now()).Scan(&cartID)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

// AddItem adds or increments quantity for product in user's cart
func (s *CartService) AddItem(userID int, productID int64, quantity int) error {
	ctx := context.Background()
	if quantity <= 0 {
		return errors.New("quantity must be > 0")
	}
	cartID, err := s.ensureCart(ctx, userID)
	if err != nil {
		return err
	}

	// upsert: if unique(cart_id,product_id) exists, increment quantity
	_, err = s.DB.Exec(ctx, `
	INSERT INTO cart_items (cart_id, product_id, quantity, created_at, updated_at)
	VALUES ($1,$2,$3,$4,$5)
	ON CONFLICT (cart_id, product_id)
	DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = EXCLUDED.updated_at
	`, cartID, productID, quantity, time.Now(), time.Now())
	return err
}

func (s *CartService) GetCartItems(userID int) ([]models.CartItem, error) {
	ctx := context.Background()
	var cartID int
	if err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID); err != nil {
		// no cart -> empty slice
		return []models.CartItem{}, nil
	}

	rows, err := s.DB.Query(ctx, "SELECT id, cart_id, product_id, quantity, created_at, updated_at FROM cart_items WHERE cart_id=$1", cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.CartItem{}
	for rows.Next() {
		var it models.CartItem
		if err := rows.Scan(&it.ID, &it.CartID, &it.ProductID, &it.Quantity, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, nil
}

func (s *CartService) RemoveItem(itemID int) error {
	ctx := context.Background()
	_, err := s.DB.Exec(ctx, "DELETE FROM cart_items WHERE id=$1", itemID)
	return err
}
