package user

import (
	"context"
	"fmt"

	"github.com/MitrofanK/Test-Menti/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type Token interface {
	GenerateToken(userID int, roleID int) (string, error)
}

type UserService struct {
	repo  Repository
	token Token
}

func NewService(repo Repository, token Token) *UserService {
	return &UserService{
		repo:  repo,
		token: token,
	}
}

func (s *UserService) SignUp(ctx context.Context, email, password string) (int, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return 0, fmt.Errorf("user with email '%s' already exists", email)
	}

	if password == "" {
		return 0, fmt.Errorf("the password field cannot be empty")
	}

	if len(password) < 8 || len(password) > 64 {
		return 0, fmt.Errorf("password must be between 8 and 64 characters long")
	}

	if password == email {
		return 0, fmt.Errorf("login and password must not match")
	}

	// проверка на состав пороля

	// проверка по топ-1000 популярных поролей (геморой)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	var newUser = models.User{
		Email:        email,
		PasswordHash: string(hashedPassword),
		RoleID:       2,
	}

	id, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}

func (s *UserService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("incorrect login or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("incorrect login or password")
	}

	token, err := s.token.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return "", fmt.Errorf("token generation error: %w", err)
	}

	return token, nil
}
