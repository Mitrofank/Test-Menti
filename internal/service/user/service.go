package user

import (
	"context"
	"fmt"

	"github.com/MitrofanK/Test-Menti/internal/errorsx"
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
		return 0, errorsx.ErrUserExists
	}

	if password == email {
		return 0, errorsx.ErrPasAndLoginSame
	}

	// проверка на состав пороля

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
		return "", errorsx.ErrIncorLoginOrPas
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errorsx.ErrIncorLoginOrPas
	}

	token, err := s.token.GenerateToken(user.ID, user.RoleID)
	if err != nil {
		return "", fmt.Errorf("token generation error: %w", err)
	}

	return token, nil
}
