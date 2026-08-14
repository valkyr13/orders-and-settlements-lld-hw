package payments

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Order struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TotalCents int64
}

type Payment struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	AmountCents int64     `json:"amount_cents"`
	PaymentDate time.Time `json:"payment_date"`
	Note        *string   `json:"note,omitempty"`
}

type Repository interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)

	GetOrderForUpdate(
		ctx context.Context,
		tx pgx.Tx,
		orderID, userID uuid.UUID,
	) (*Order, error)

	GetAmountPaid(
		ctx context.Context,
		tx pgx.Tx,
		orderID uuid.UUID,
	) (int64, error)

	CreatePayment(
		ctx context.Context,
		tx pgx.Tx,
		payment Payment,
	) error

	ListPayments(
		ctx context.Context,
		orderID, userID uuid.UUID,
	) ([]Payment, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

func (r *repository) GetOrderForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	orderID, userID uuid.UUID,
) (*Order, error) {
	var order Order

	err := tx.QueryRow(
		ctx,
		`SELECT id, user_id, total_cents
		 FROM orders
		 WHERE id = $1 AND user_id = $2
		 FOR UPDATE`,
		orderID,
		userID,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.TotalCents,
	)

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *repository) GetAmountPaid(
	ctx context.Context,
	tx pgx.Tx,
	orderID uuid.UUID,
) (int64, error) {
	var amountPaid int64

	err := tx.QueryRow(
		ctx,
		`SELECT COALESCE(SUM(amount_cents), 0)
		 FROM payments
		 WHERE order_id = $1`,
		orderID,
	).Scan(&amountPaid)

	return amountPaid, err
}

func (r *repository) CreatePayment(
	ctx context.Context,
	tx pgx.Tx,
	payment Payment,
) error {
	_, err := tx.Exec(
		ctx,
		`INSERT INTO payments
			(id, order_id, amount_cents, payment_date, note)
		 VALUES ($1, $2, $3, $4, $5)`,
		payment.ID,
		payment.OrderID,
		payment.AmountCents,
		payment.PaymentDate,
		payment.Note,
	)

	return err
}

func (r *repository) ListPayments(
	ctx context.Context,
	orderID, userID uuid.UUID,
) ([]Payment, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
			p.id,
			p.order_id,
			p.amount_cents,
			p.payment_date,
			p.note
		 FROM payments p
		 JOIN orders o ON o.id = p.order_id
		 WHERE p.order_id = $1
		   AND o.user_id = $2
		 ORDER BY p.payment_date DESC`,
		orderID,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []Payment

	for rows.Next() {
		var payment Payment

		if err := rows.Scan(
			&payment.ID,
			&payment.OrderID,
			&payment.AmountCents,
			&payment.PaymentDate,
			&payment.Note,
		); err != nil {
			return nil, err
		}

		payments = append(payments, payment)
	}

	return payments, rows.Err()
}