package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type ShiftRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Shift, error)
	FindOpenByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error)
	FindAll(ctx context.Context, cashierID string, page, limit int) ([]entity.Shift, int, error)
	Create(ctx context.Context, shift *entity.Shift) error
	Close(ctx context.Context, id string, closingCash int64, notes string) error
}

type shiftRepository struct {
	db *sql.DB
}

func NewShiftRepository(db *sql.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) FindByID(ctx context.Context, id string) (*entity.Shift, error) {
	shift := &entity.Shift{}
	cashier := &entity.User{}

	query := `
		SELECT s.id, s.cashier_id, s.opened_at, s.closed_at, s.opening_cash, s.closing_cash, 
		       s.total_transactions, s.status, s.notes, s.created_at, s.updated_at,
		       u.id, u.name, u.email, u.role, u.is_active
		FROM shifts s
		JOIN users u ON u.id = s.cashier_id
		WHERE s.id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shift.ID, &shift.CashierID, &shift.OpenedAt, &shift.ClosedAt, &shift.OpeningCash, &shift.ClosingCash,
		&shift.TotalTransactions, &shift.Status, &shift.Notes, &shift.CreatedAt, &shift.UpdatedAt,
		&cashier.ID, &cashier.Name, &cashier.Email, &cashier.Role, &cashier.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	shift.Cashier = cashier
	return shift, nil
}

func (r *shiftRepository) FindOpenByCashierID(ctx context.Context, cashierID string) (*entity.Shift, error) {
	shift := &entity.Shift{}

	query := `
		SELECT id, cashier_id, opened_at, closed_at, opening_cash, closing_cash, 
		       total_transactions, status, notes, created_at, updated_at
		FROM shifts 
		WHERE cashier_id = ? AND status = 'open'
	`

	err := r.db.QueryRowContext(ctx, query, cashierID).Scan(
		&shift.ID, &shift.CashierID, &shift.OpenedAt, &shift.ClosedAt, &shift.OpeningCash, &shift.ClosingCash,
		&shift.TotalTransactions, &shift.Status, &shift.Notes, &shift.CreatedAt, &shift.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return shift, nil
}

func (r *shiftRepository) FindAll(ctx context.Context, cashierID string, page, limit int) ([]entity.Shift, int, error) {
	var whereClauses []string
	var args []interface{}

	if cashierID != "" {
		whereClauses = append(whereClauses, "s.cashier_id = ?")
		args = append(args, cashierID)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total records
	countQuery := `SELECT COUNT(*) FROM shifts s ` + whereSQL
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Pagination setup
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	selectSQL := `
		SELECT s.id, s.cashier_id, s.opened_at, s.closed_at, s.opening_cash, s.closing_cash, 
		       s.total_transactions, s.status, s.notes, s.created_at, s.updated_at,
		       u.id, u.name, u.email, u.role, u.is_active
		FROM shifts s
		JOIN users u ON u.id = s.cashier_id
	` + whereSQL + ` ORDER BY s.opened_at DESC LIMIT ? OFFSET ?`

	selectArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	shifts := []entity.Shift{}
	for rows.Next() {
		var s entity.Shift
		var u entity.User
		err := rows.Scan(
			&s.ID, &s.CashierID, &s.OpenedAt, &s.ClosedAt, &s.OpeningCash, &s.ClosingCash,
			&s.TotalTransactions, &s.Status, &s.Notes, &s.CreatedAt, &s.UpdatedAt,
			&u.ID, &u.Name, &u.Email, &u.Role, &u.IsActive,
		)
		if err != nil {
			return nil, 0, err
		}
		s.Cashier = &u
		shifts = append(shifts, s)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return shifts, total, nil
}

func (r *shiftRepository) Create(ctx context.Context, shift *entity.Shift) error {
	shift.ID = uuid.New().String()
	shift.OpenedAt = time.Now()
	shift.CreatedAt = time.Now()
	shift.UpdatedAt = time.Now()
	shift.Status = entity.ShiftStatusOpen
	shift.TotalTransactions = 0

	query := `
		INSERT INTO shifts (
			id, cashier_id, opened_at, opening_cash, total_transactions, status, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		shift.ID, shift.CashierID, shift.OpenedAt, shift.OpeningCash, shift.TotalTransactions, shift.Status, shift.Notes, shift.CreatedAt, shift.UpdatedAt,
	)
	return err
}

func (r *shiftRepository) Close(ctx context.Context, id string, closingCash int64, notes string) error {
	query := `
		UPDATE shifts SET 
			status = 'closed',
			closed_at = NOW(),
			closing_cash = ?,
			notes = ?,
			updated_at = NOW()
		WHERE id = ? AND status = 'open'
	`

	res, err := r.db.ExecContext(ctx, query, closingCash, notes, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("shift not found or already closed")
	}

	return nil
}
