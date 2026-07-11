package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type PaymentRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error)
	FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*entity.Payment, error)
	Create(ctx context.Context, payment *entity.Payment) error
	UpdateStatus(ctx context.Context, id string, status string, paidAt *time.Time) error
	UpdateMidtransData(ctx context.Context, id string, token string, url string, midtransOrderID string) error
	SaveRawNotification(ctx context.Context, id string, raw string) error
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) FindByID(ctx context.Context, id string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	query := `
		SELECT id, order_id, method, status, amount, midtrans_order_id, 
		       midtrans_token, midtrans_url, raw_notification, paid_at, created_at, updated_at
		FROM payments
		WHERE id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID, &payment.OrderID, &payment.Method, &payment.Status, &payment.Amount,
		&payment.MidtransOrderID, &payment.MidtransToken, &payment.MidtransURL,
		&payment.RawNotification, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	query := `
		SELECT id, order_id, method, status, amount, midtrans_order_id, 
		       midtrans_token, midtrans_url, raw_notification, paid_at, created_at, updated_at
		FROM payments
		WHERE order_id = ?
	`

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID, &payment.OrderID, &payment.Method, &payment.Status, &payment.Amount,
		&payment.MidtransOrderID, &payment.MidtransToken, &payment.MidtransURL,
		&payment.RawNotification, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) FindByMidtransOrderID(ctx context.Context, midtransOrderID string) (*entity.Payment, error) {
	payment := &entity.Payment{}
	query := `
		SELECT id, order_id, method, status, amount, midtrans_order_id, 
		       midtrans_token, midtrans_url, raw_notification, paid_at, created_at, updated_at
		FROM payments
		WHERE midtrans_order_id = ?
	`

	err := r.db.QueryRowContext(ctx, query, midtransOrderID).Scan(
		&payment.ID, &payment.OrderID, &payment.Method, &payment.Status, &payment.Amount,
		&payment.MidtransOrderID, &payment.MidtransToken, &payment.MidtransURL,
		&payment.RawNotification, &payment.PaidAt, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) Create(ctx context.Context, payment *entity.Payment) error {
	payment.ID = uuid.New().String()
	payment.CreatedAt = time.Now()
	payment.UpdatedAt = time.Now()
	if payment.Status == "" {
		payment.Status = entity.PaymentStatusPending
	}

	query := `
		INSERT INTO payments (
			id, order_id, method, status, amount, midtrans_order_id, 
			midtrans_token, midtrans_url, raw_notification, paid_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		payment.ID, payment.OrderID, payment.Method, payment.Status, payment.Amount,
		payment.MidtransOrderID, payment.MidtransToken, payment.MidtransURL,
		payment.RawNotification, payment.PaidAt, payment.CreatedAt, payment.UpdatedAt,
	)
	return err
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, status string, paidAt *time.Time) error {
	query := `
		UPDATE payments 
		SET status = ?, paid_at = ?, updated_at = NOW() 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, status, paidAt, id)
	return err
}

func (r *paymentRepository) UpdateMidtransData(ctx context.Context, id string, token string, url string, midtransOrderID string) error {
	query := `
		UPDATE payments 
		SET midtrans_token = ?, midtrans_url = ?, midtrans_order_id = ?, updated_at = NOW() 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, token, url, midtransOrderID, id)
	return err
}

func (r *paymentRepository) SaveRawNotification(ctx context.Context, id string, raw string) error {
	query := `
		UPDATE payments 
		SET raw_notification = ?, updated_at = NOW() 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, raw, id)
	return err
}
