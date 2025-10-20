package services

import (
	"context"
	"fmt"
	"time"

	"github.com/edsjcbra/flightsimhub/config"
	"github.com/edsjcbra/flightsimhub/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{DB: db}
}

func (s *AuthService) CreateUser(name, email, password string) (*models.User, error) {
	ctx := context.Background()

	var exists int
	if err := s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&exists); err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("email already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var id int
	if err := s.DB.QueryRow(ctx,
		"INSERT INTO users (name,email,password_hash,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		name, email, string(hashed), time.Now(), time.Now(),
	).Scan(&id); err != nil {
		return nil, err
	}

	u := &models.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  string(hashed),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return u, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	ctx := context.Background()

	var id int
	var hashed string
	if err := s.DB.QueryRow(ctx, "SELECT id, password_hash FROM users WHERE email=$1", email).Scan(&id, &hashed); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
