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

type UserService struct {
	repo Repository
}

func NewService(repo Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) SignUp(ctx context.Context, email, pasword string) (int, error) {
	if _, err := s.repo.GetUserByEmail(ctx, email); err == nil {
		return 0, fmt.Errorf("user with email '%s' already exists", email)
	}

	if pasword == "" {
		return 0, fmt.Errorf("the password field cannot be empty")
	}

	if len(pasword) < 8 || len(pasword) > 64 {
		return 0, fmt.Errorf("password must be between 8 and 64 characters long")
	}

	if pasword == email {
		return 0, fmt.Errorf("login and password must not match")
	}

	// проверка на состав пороля

	// проверка по топ-1000 популярных поролей (геморой)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pasword), bcrypt.DefaultCost)
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
