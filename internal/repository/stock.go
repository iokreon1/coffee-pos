package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type StockFilter struct {
	ProductID string
	UserID    string
	Type      string
	Page      int
	Limit     int
}

type StockRepository interface {
	Create(ctx context.Context, movement *entity.StockMovement) error
	FindByProductID(ctx context.Context, productID string, filter StockFilter) ([]entity.StockMovement, int, error)
	FindAll(ctx context.Context, filter StockFilter) ([]entity.StockMovement, int, error)
}

type stockRepository struct {
	db *sql.DB
}

func NewStockRepository(db *sql.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(ctx context.Context, movement *entity.StockMovement) error {
	movement.ID = uuid.New().String()
	movement.CreatedAt = time.Now()
	movement.UpdatedAt = time.Now()

	query := `
		INSERT INTO stock_movements (
			id, product_id, user_id, type, quantity, 
			stock_before, stock_after, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		movement.ID, movement.ProductID, movement.UserID, movement.Type, movement.Quantity,
		movement.StockBefore, movement.StockAfter, movement.Notes, movement.CreatedAt, movement.UpdatedAt,
	)
	return err
}

func (r *stockRepository) FindByProductID(ctx context.Context, productID string, filter StockFilter) ([]entity.StockMovement, int, error) {
	filter.ProductID = productID
	return r.FindAll(ctx, filter)
}

func (r *stockRepository) FindAll(ctx context.Context, filter StockFilter) ([]entity.StockMovement, int, error) {
	var whereClauses []string
	var args []interface{}

	if filter.ProductID != "" {
		whereClauses = append(whereClauses, "sm.product_id = ?")
		args = append(args, filter.ProductID)
	}
	if filter.UserID != "" {
		whereClauses = append(whereClauses, "sm.user_id = ?")
		args = append(args, filter.UserID)
	}
	if filter.Type != "" {
		whereClauses = append(whereClauses, "sm.type = ?")
		args = append(args, filter.Type)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total movements
	countQuery := `SELECT COUNT(*) FROM stock_movements sm ` + whereSQL
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Pagination setup
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	selectSQL := `
		SELECT sm.id, sm.product_id, sm.user_id, sm.type, sm.quantity,
		       sm.stock_before, sm.stock_after, sm.notes, sm.created_at, sm.updated_at,
		       p.id, p.name,
		       u.id, u.name
		FROM stock_movements sm
		JOIN products p ON p.id = sm.product_id
		JOIN users u ON u.id = sm.user_id
	` + whereSQL + ` ORDER BY sm.created_at DESC LIMIT ? OFFSET ?`

	selectArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	movements := []entity.StockMovement{}
	for rows.Next() {
		var sm entity.StockMovement
		var p entity.Product
		var u entity.User

		err := rows.Scan(
			&sm.ID, &sm.ProductID, &sm.UserID, &sm.Type, &sm.Quantity,
			&sm.StockBefore, &sm.StockAfter, &sm.Notes, &sm.CreatedAt, &sm.UpdatedAt,
			&p.ID, &p.Name,
			&u.ID, &u.Name,
		)
		if err != nil {
			return nil, 0, err
		}

		sm.Product = &p
		sm.User = &u
		movements = append(movements, sm)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return movements, total, nil
}
