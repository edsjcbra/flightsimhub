package models

import "time"

type Cart struct {
    ID        int       `json:"id"`
    UserID    int       `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
    ID        int       `json:"id"`
    CartID    int       `json:"cart_id"`
    ProductID int       `json:"product_id"`
    Quantity  int       `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
