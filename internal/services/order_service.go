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

// Novo construtor usando *pgxpool.Pool direto
func NewOrderService(db *pgxpool.Pool) *OrderService {
	return &OrderService{DB: db}
}

// CreateOrder converte o carrinho do usuário em pedido
func (s *OrderService) CreateOrder(userID int) (models.Order, error) {
	ctx := context.Background()

	// Busca carrinho do usuário
	var cartID int
	err := s.DB.QueryRow(ctx, "SELECT id FROM carts WHERE user_id=$1", userID).Scan(&cartID)
	if err != nil {
		return models.Order{}, err
	}

	rows, err := s.DB.Query(ctx, "SELECT product_id, quantity FROM cart_items WHERE cart_id=$1", cartID)
	if err != nil {
		return models.Order{}, err
	}
	defer rows.Close()

	var total float64
	var items []models.OrderItem
	for rows.Next() {
		var productID, quantity int
		if err = rows.Scan(&productID, &quantity); err != nil {
			return models.Order{}, err
		}

		var price float64
		err = s.DB.QueryRow(ctx, "SELECT price FROM products WHERE id=$1", productID).Scan(&price)
		if err != nil {
			return models.Order{}, err
		}

		total += price * float64(quantity)
		items = append(items, models.OrderItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     price,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	// Cria pedido
	var orderID int
	err = s.DB.QueryRow(ctx, "INSERT INTO orders (user_id, total, created_at, updated_at) VALUES ($1,$2,$3,$4) RETURNING id",
		userID, total, time.Now(), time.Now()).Scan(&orderID)
	if err != nil {
		return models.Order{}, err
	}

	// Insere itens do pedido
	for _, item := range items {
		_, err := s.DB.Exec(ctx,
			"INSERT INTO order_items (order_id, product_id, quantity, price, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6)",
			orderID, item.ProductID, item.Quantity, item.Price, item.CreatedAt, item.UpdatedAt)
		if err != nil {
			return models.Order{}, err
		}
	}

	// Limpa carrinho
	_, err = s.DB.Exec(ctx, "DELETE FROM cart_items WHERE cart_id=$1", cartID)
	if err != nil {
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

// GetOrders retorna todos os pedidos do usuário
func (s *OrderService) GetOrders(userID int) ([]models.Order, error) {
	ctx := context.Background()
	rows, err := s.DB.Query(ctx, "SELECT id, user_id, total, created_at, updated_at FROM orders WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.UserID, &order.Total, &order.CreatedAt, &order.UpdatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}
