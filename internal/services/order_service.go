package services

import (
	"context"
	"time"

	"github.com/edsjcbra/flightsimhub/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	DB *pgxpool.Pool
}

func NewOrderService(db *pgxpool.Pool) *OrderService {
	return &OrderService{DB: db}
}

// CreateOrder converts user's cart into an order, returns Order
func (s *OrderService) CreateOrder(userID int) (models.Order, error) {
	ctx := context.Background()

	// get cart id
	var cartID int
	if err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID); err != nil {
		return models.Order{}, err
	}

	// read cart items
	rows, err := s.DB.Query(ctx, "SELECT product_id, quantity FROM cart_items WHERE cart_id=$1", cartID)
	if err != nil {
		return models.Order{}, err
	}
	defer rows.Close()

	var items []models.OrderItem
	var total float64
	for rows.Next() {
		var pid int64
		var qty int
		if err := rows.Scan(&pid, &qty); err != nil {
			return models.Order{}, err
		}
		// get price
		var price float64
		if err := s.DB.QueryRow(ctx, "SELECT price FROM products WHERE id=$1", pid).Scan(&price); err != nil {
			return models.Order{}, err
		}
		total += price * float64(qty)
		items = append(items, models.OrderItem{
			ProductID: pid,
			Quantity:  qty,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	// create order
	var orderID int
	if err := s.DB.QueryRow(ctx, "INSERT INTO orders (user_id,total,created_at,updated_at) VALUES ($1,$2,$3,$4) RETURNING id",
		userID, total, time.Now(), time.Now()).Scan(&orderID); err != nil {
		return models.Order{}, err
	}

	// insert items
	for _, it := range items {
		if _, err := s.DB.Exec(ctx, "INSERT INTO order_items (order_id,product_id,quantity,price,created_at,updated_at) VALUES ($1,$2,$3,$4,$5,$6)",
			orderID, it.ProductID, it.Quantity, it.Price, it.CreatedAt, it.UpdatedAt); err != nil {
			return models.Order{}, err
		}
	}

	// clear cart items
	if _, err := s.DB.Exec(ctx, "DELETE FROM cart_items WHERE cart_id=$1", cartID); err != nil {
		return models.Order{}, err
	}

	return models.Order{
		ID:        orderID,
		UserID:    userID,
		Total:     total,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	ctx := context.Background()
	rows, err := s.DB.Query(ctx, "SELECT id,user_id,total,created_at,updated_at FROM orders WHERE user_id=$1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.UserID, &o.Total, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	return orders, nil
}
