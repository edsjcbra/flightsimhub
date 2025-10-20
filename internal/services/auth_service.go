package services

import (
	"context"
	"errors"
	"time"

	"github.com/edsjcbra/flightsimhub/config"
	"github.com/edsjcbra/flightsimhub/internal/database"
	"github.com/edsjcbra/flightsimhub/internal/models"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Criar novo usuário
func (s *AuthService) CreateUser(name, email, password string) (*models.User, error) {
	ctx := context.Background()

	// Verificar se email já existe
	var exists int
	err := database.DB.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE email=$1", email).Scan(&exists)
	if err != nil {
		return nil, err
	}
	if exists > 0 {
		return nil, errors.New("email already registered")
	}

	// Hashear senha
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	var id int64
	err = database.DB.QueryRow(ctx,
		"INSERT INTO users (name, email, password_hash, created_at, updated_at) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		name, email, string(hash), now, now).Scan(&id)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	return user, nil
}

// Autenticar usuário e gerar JWT
func (s *AuthService) AuthenticateUser(email, password string) (string, error) {
	ctx := context.Background()
	var user models.User
	err := database.DB.QueryRow(ctx, "SELECT id, password_hash FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.PasswordHash)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Criar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // 3 dias
	})

	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
