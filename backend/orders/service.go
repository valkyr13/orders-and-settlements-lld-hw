package orders

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidOrder = errors.New("invalid order")
)

const (
	StatusPending       = "PENDING"
	StatusPartiallyPaid = "PARTIALLY_PAID"
	StatusOverdue       = "OVERDUE"
	StatusPaid          = "PAID"
)

type LineItem struct {
	ID             uuid.UUID `json:"id"`
	Description    string    `json:"description"`
	Quantity       int       `json:"quantity"`
	UnitPriceCents int64     `json:"unit_price_cents"`
}

type Order struct {
	ID              uuid.UUID  `json:"id"`
	CustomerName    string     `json:"customer_name"`
	DueDate         time.Time  `json:"due_date"`
	TotalCents      int64      `json:"total_cents"`
	AmountPaidCents int64      `json:"amount_paid_cents"`
	AmountDueCents  int64      `json:"amount_due_cents"`
	Status          string     `json:"status"`
	LineItems       []LineItem `json:"line_items"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateOrder(
	ctx context.Context,
	userID uuid.UUID,
	customerName string,
	dueDate time.Time,
	items []LineItem,
) error {
	if customerName == "" || len(items) == 0 {
		return ErrInvalidOrder
	}

	var totalCents int64

	for _, item := range items {
		if item.Quantity < 1 || item.UnitPriceCents < 0 {
			return ErrInvalidOrder
		}

		totalCents += int64(item.Quantity) * item.UnitPriceCents
	}

	if totalCents <= 0 {
		return ErrInvalidOrder
	}

	return s.repo.CreateOrder(
		ctx,
		uuid.New(),
		userID,
		customerName,
		dueDate.Format(time.RFC3339),
		totalCents,
		items,
	)
}

func (s *Service) ListOrders(
	ctx context.Context,
	userID uuid.UUID,
) ([]Order, error) {
	return s.repo.ListOrders(ctx, userID)
}

func (s *Service) GetOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
) (*Order, error) {
	return s.repo.GetOrder(ctx, userID, orderID)
}

func (s *Service) UpdateOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
	customerName string,
	dueDate time.Time,
	items []LineItem,
) error {
	totalCents, err := calculateTotal(items)
	if err != nil {
		return err
	}

	return s.repo.UpdateOrder(
		ctx,
		userID,
		orderID,
		customerName,
		dueDate,
		totalCents,
		items,
	)
}

func (s *Service) DeleteOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
) error {
	return s.repo.DeleteOrder(ctx, userID, orderID)
}

func calculateTotal(items []LineItem) (int64, error) {
	if len(items) == 0 {
		return 0, ErrInvalidOrder
	}

	var total int64

	for _, item := range items {
		if item.Quantity < 1 || item.UnitPriceCents < 0 {
			return 0, ErrInvalidOrder
		}

		total += int64(item.Quantity) * item.UnitPriceCents
	}

	if total <= 0 {
		return 0, ErrInvalidOrder
	}

	return total, nil
}

func deriveStatus(
	totalCents int64,
	amountPaidCents int64,
	dueDate time.Time,
	now time.Time,
) string {
	if amountPaidCents >= totalCents {
		return StatusPaid
	}

	if dueDate.Before(now) {
		return StatusOverdue
	}

	if amountPaidCents > 0 {
		return StatusPartiallyPaid
	}

	return StatusPending
}
