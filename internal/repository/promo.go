package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type PromoRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Promo, error)
	FindByCode(ctx context.Context, code string) (*entity.Promo, error)
	FindAll(ctx context.Context, page, limit int) ([]entity.Promo, int, error)
	Create(ctx context.Context, promo *entity.Promo) error
	Update(ctx context.Context, promo *entity.Promo) error
	Delete(ctx context.Context, id string) error
	IncrementUsedCount(ctx context.Context, id string) error
}

type promoRepository struct {
	db *sql.DB
}

func NewPromoRepository(db *sql.DB) PromoRepository {
	return &promoRepository{db: db}
}

func (r *promoRepository) FindByID(ctx context.Context, id string) (*entity.Promo, error) {
	promo := &entity.Promo{}
	query := `
		SELECT id, name, code, type, value, min_order, max_discount, usage_limit, used_count, 
		       started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos 
		WHERE id = ? AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&promo.ID, &promo.Name, &promo.Code, &promo.Type, &promo.Value, &promo.MinOrder,
		&promo.MaxDiscount, &promo.UsageLimit, &promo.UsedCount,
		&promo.StartedAt, &promo.EndedAt, &promo.IsActive, &promo.CreatedAt, &promo.UpdatedAt, &promo.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return promo, nil
}

func (r *promoRepository) FindByCode(ctx context.Context, code string) (*entity.Promo, error) {
	promo := &entity.Promo{}
	query := `
		SELECT id, name, code, type, value, min_order, max_discount, usage_limit, used_count, 
		       started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos 
		WHERE code = ? AND deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&promo.ID, &promo.Name, &promo.Code, &promo.Type, &promo.Value, &promo.MinOrder,
		&promo.MaxDiscount, &promo.UsageLimit, &promo.UsedCount,
		&promo.StartedAt, &promo.EndedAt, &promo.IsActive, &promo.CreatedAt, &promo.UpdatedAt, &promo.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return promo, nil
}

func (r *promoRepository) FindAll(ctx context.Context, page, limit int) ([]entity.Promo, int, error) {
	// Count total records
	countQuery := `SELECT COUNT(*) FROM promos WHERE deleted_at IS NULL`
	var total int
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
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
		SELECT id, name, code, type, value, min_order, max_discount, usage_limit, used_count, 
		       started_at, ended_at, is_active, created_at, updated_at, deleted_at
		FROM promos 
		WHERE deleted_at IS NULL 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, selectSQL, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	promos := []entity.Promo{}
	for rows.Next() {
		var promo entity.Promo
		err := rows.Scan(
			&promo.ID, &promo.Name, &promo.Code, &promo.Type, &promo.Value, &promo.MinOrder,
			&promo.MaxDiscount, &promo.UsageLimit, &promo.UsedCount,
			&promo.StartedAt, &promo.EndedAt, &promo.IsActive, &promo.CreatedAt, &promo.UpdatedAt, &promo.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		promos = append(promos, promo)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return promos, total, nil
}

func (r *promoRepository) Create(ctx context.Context, promo *entity.Promo) error {
	promo.ID = uuid.New().String()
	promo.CreatedAt = time.Now()
	promo.UpdatedAt = time.Now()
	promo.UsedCount = 0

	query := `
		INSERT INTO promos (
			id, name, code, type, value, min_order, max_discount, usage_limit, used_count, 
			started_at, ended_at, is_active, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		promo.ID, promo.Name, promo.Code, promo.Type, promo.Value, promo.MinOrder,
		promo.MaxDiscount, promo.UsageLimit, promo.UsedCount,
		promo.StartedAt, promo.EndedAt, promo.IsActive, promo.CreatedAt, promo.UpdatedAt,
	)
	return err
}

func (r *promoRepository) Update(ctx context.Context, promo *entity.Promo) error {
	promo.UpdatedAt = time.Now()

	query := `
		UPDATE promos SET 
			name = ?, type = ?, value = ?, min_order = ?, max_discount = ?, usage_limit = ?, 
			started_at = ?, ended_at = ?, is_active = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		promo.Name, promo.Type, promo.Value, promo.MinOrder, promo.MaxDiscount, promo.UsageLimit,
		promo.StartedAt, promo.EndedAt, promo.IsActive, promo.UpdatedAt, promo.ID,
	)
	return err
}

func (r *promoRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE promos SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *promoRepository) IncrementUsedCount(ctx context.Context, id string) error {
	query := `
		UPDATE promos SET used_count = used_count + 1, updated_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
