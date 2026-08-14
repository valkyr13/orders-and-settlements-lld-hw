package payments

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	ErrOrderNotFound     = errors.New("order not found")
	ErrInvalidPayment    = errors.New("invalid payment")
	ErrPaymentExceedsDue = errors.New("payment exceeds amount due")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePayment(
	ctx context.Context,
	userID, orderID uuid.UUID,
	amountCents int64,
	paymentDate time.Time,
	note *string,
) (*Payment, error) {
	if amountCents <= 0 {
		return nil, ErrInvalidPayment
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Locks only this order row.
	order, err := s.repo.GetOrderForUpdate(
		ctx,
		tx,
		orderID,
		userID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrderNotFound
		}

		return nil, err
	}

	amountPaid, err := s.repo.GetAmountPaid(
		ctx,
		tx,
		orderID,
	)
	if err != nil {
		return nil, err
	}

	amountDue := order.TotalCents - amountPaid

	if amountCents > amountDue {
		return nil, ErrPaymentExceedsDue
	}

	payment := Payment{
		ID:          uuid.New(),
		OrderID:     orderID,
		AmountCents: amountCents,
		PaymentDate: paymentDate,
		Note:        note,
	}

	if err := s.repo.CreatePayment(ctx, tx, payment); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (s *Service) ListPayments(
	ctx context.Context,
	userID, orderID uuid.UUID,
) ([]Payment, error) {
	return s.repo.ListPayments(ctx, orderID, userID)
}
