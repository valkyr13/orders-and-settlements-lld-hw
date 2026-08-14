package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Signup(
	ctx context.Context,
	email, password string,
) error {
	exists, err := s.repo.UserExists(ctx, email)
	if err != nil {
		return err
	}

	if exists {
		return ErrUserExists
	}

	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}

	return s.repo.CreateUser(
		ctx,
		uuid.New(),
		email,
		string(passwordHash),
	)
}

func (s *Service) Login(
	ctx context.Context,
	email, password string,
) (uuid.UUID, error) {
	userID, passwordHash, err := s.repo.GetUserCredentials(
		ctx,
		email,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, ErrInvalidCredentials
	}

	if err != nil {
		return uuid.Nil, err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		[]byte(password),
	); err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	return userID, nil
}

func (s *Service) CreateSession(
	ctx context.Context,
	userID uuid.UUID,
) (uuid.UUID, error) {
	sessionID := uuid.New()

	err := s.repo.CreateSession(
		ctx,
		sessionID,
		userID,
	)

	if err != nil {
		return uuid.Nil, err
	}

	return sessionID, nil
}
