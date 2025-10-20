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

// CreateUser cria usuário
func (s *AuthService) CreateUser(name, email, password string) (*models.User, error) {
	ctx := context.Background()

	var exists int
	err := s.DB.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var id int
	err = s.DB.QueryRow(ctx,
		"INSERT INTO users (name,email,password,created_at,updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		name, email, string(hashedPassword), time.Now(), time.Now(),
	).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// Login retorna JWT
func (s *AuthService) Login(email, password string) (string, error) {
	ctx := context.Background()

	var user models.User
	err := s.DB.QueryRow(ctx, "SELECT id,password FROM users WHERE email=$1", email).Scan(&user.ID, &user.Password)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}
