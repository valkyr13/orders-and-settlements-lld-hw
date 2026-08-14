package orders

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrOrderHasPayments = errors.New("order has payments")
)

type Repository interface {
	CreateOrder(
		ctx context.Context,
		orderID, userID uuid.UUID,
		customerName string,
		dueDate string,
		totalCents int64,
		items []LineItem,
	) error

	ListOrders(
		ctx context.Context,
		userID uuid.UUID,
	) ([]Order, error)

	GetOrder(
		ctx context.Context,
		userID, orderID uuid.UUID,
	) (*Order, error)
	UpdateOrder(
		ctx context.Context,
		userID, orderID uuid.UUID,
		customerName string,
		dueDate time.Time,
		totalCents int64,
		items []LineItem,
	) error
	DeleteOrder(ctx context.Context, userID, orderID uuid.UUID) error
	HasPayments(ctx context.Context, orderID uuid.UUID) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrder(
	ctx context.Context,
	orderID, userID uuid.UUID,
	customerName string,
	dueDate string,
	totalCents int64,
	items []LineItem,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		`INSERT INTO orders
			(id, user_id, customer_name, due_date, total_cents)
		 VALUES ($1, $2, $3, $4, $5)`,
		orderID,
		userID,
		customerName,
		dueDate,
		totalCents,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO order_items
				(id, order_id, description, quantity, unit_price_cents)
			 VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(),
			orderID,
			item.Description,
			item.Quantity,
			item.UnitPriceCents,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) ListOrders(
	ctx context.Context,
	userID uuid.UUID,
) ([]Order, error) {

	query := `SELECT
    o.id,
    o.customer_name,
    o.due_date,
    o.total_cents,
    COALESCE(SUM(p.amount_cents), 0) AS amount_paid_cents,
    o.created_at,
    o.updated_at
FROM orders o
LEFT JOIN payments p ON p.order_id = o.id
WHERE o.user_id = $1
GROUP BY
    o.id,
    o.customer_name,
    o.due_date,
    o.total_cents,
    o.created_at,
    o.updated_at
ORDER BY o.created_at DESC`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var order Order

		if err := rows.Scan(
			&order.ID,
			&order.CustomerName,
			&order.DueDate,
			&order.TotalCents,
			&order.AmountPaidCents,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}

		order.AmountDueCents = order.TotalCents - order.AmountPaidCents
		order.Status = deriveStatus(
			order.TotalCents,
			order.AmountPaidCents,
			order.DueDate,
			time.Now(),
		)

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *repository) GetOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
) (*Order, error) {
	var order Order

	query := `SELECT
    o.id,
    o.customer_name,
    o.due_date,
    o.total_cents,
    COALESCE(SUM(p.amount_cents), 0) AS amount_paid_cents,
    o.created_at,
    o.updated_at
FROM orders o
LEFT JOIN payments p ON p.order_id = o.id
WHERE o.id = $1
  AND o.user_id = $2
GROUP BY
    o.id,
    o.customer_name,
    o.due_date,
    o.total_cents,
    o.created_at,
    o.updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		orderID,
		userID,
	).Scan(
		&order.ID,
		&order.CustomerName,
		&order.DueDate,
		&order.TotalCents,
		&order.AmountPaidCents,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(
		ctx,
		`SELECT
			id,
			description,
			quantity,
			unit_price_cents
		 FROM order_items
		 WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item LineItem

		if err := rows.Scan(
			&item.ID,
			&item.Description,
			&item.Quantity,
			&item.UnitPriceCents,
		); err != nil {
			return nil, err
		}

		order.LineItems = append(order.LineItems, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	order.AmountDueCents = order.TotalCents - order.AmountPaidCents
	order.Status = deriveStatus(
		order.TotalCents,
		order.AmountPaidCents,
		order.DueDate,
		time.Now(),
	)

	return &order, nil
}

func (r *repository) HasPayments(
	ctx context.Context,
	orderID uuid.UUID,
) (bool, error) {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1 FROM payments WHERE order_id = $1
		)`,
		orderID,
	).Scan(&exists)

	return exists, err
}

func (r *repository) UpdateOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
	customerName string,
	dueDate time.Time,
	totalCents int64,
	items []LineItem,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool

	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE id = $1 AND user_id = $2
			FOR UPDATE
		)`,
		orderID,
		userID,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return ErrOrderNotFound
	}

	var hasPayments bool

	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM payments
			WHERE order_id = $1
		)`,
		orderID,
	).Scan(&hasPayments)

	if err != nil {
		return err
	}

	if hasPayments {
		return ErrOrderHasPayments
	}

	_, err = tx.Exec(
		ctx,
		`UPDATE orders
		 SET customer_name = $1,
		     due_date = $2,
		     total_cents = $3,
		     updated_at = NOW()
		 WHERE id = $4 AND user_id = $5`,
		customerName,
		dueDate,
		totalCents,
		orderID,
		userID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM order_items
		 WHERE order_id = $1`,
		orderID,
	)
	if err != nil {
		return err
	}

	for _, item := range items {
		_, err = tx.Exec(
			ctx,
			`INSERT INTO order_items
				(id, order_id, description, quantity, unit_price_cents)
			 VALUES ($1, $2, $3, $4, $5)`,
			uuid.New(),
			orderID,
			item.Description,
			item.Quantity,
			item.UnitPriceCents,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *repository) DeleteOrder(
	ctx context.Context,
	userID, orderID uuid.UUID,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool

	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM orders
			WHERE id = $1 AND user_id = $2
			FOR UPDATE
		)`,
		orderID,
		userID,
	).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return ErrOrderNotFound
	}

	var hasPayments bool

	err = tx.QueryRow(
		ctx,
		`SELECT EXISTS (
			SELECT 1
			FROM payments
			WHERE order_id = $1
		)`,
		orderID,
	).Scan(&hasPayments)

	if err != nil {
		return err
	}

	if hasPayments {
		return ErrOrderHasPayments
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM orders
		 WHERE id = $1 AND user_id = $2`,
		orderID,
		userID,
	)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
