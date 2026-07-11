package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type TableRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Table, error)
	FindAll(ctx context.Context) ([]entity.Table, error)
	Create(ctx context.Context, table *entity.Table) error
	Update(ctx context.Context, table *entity.Table) error
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status string) error
}

type tableRepository struct {
	db *sql.DB
}

func NewTableRepository(db *sql.DB) TableRepository {
	return &tableRepository{db: db}
}

func (r *tableRepository) FindByID(ctx context.Context, id string) (*entity.Table, error) {
	table := &entity.Table{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, capacity, status, created_at, updated_at, deleted_at 
         FROM tables WHERE id = ? AND deleted_at IS NULL`, id).
		Scan(&table.ID, &table.Name, &table.Capacity, &table.Status, &table.CreatedAt, &table.UpdatedAt, &table.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (r *tableRepository) FindAll(ctx context.Context) ([]entity.Table, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, capacity, status, created_at, updated_at, deleted_at 
         FROM tables WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tables := []entity.Table{}
	for rows.Next() {
		var table entity.Table
		err := rows.Scan(&table.ID, &table.Name, &table.Capacity, &table.Status, &table.CreatedAt, &table.UpdatedAt, &table.DeletedAt)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *tableRepository) Create(ctx context.Context, table *entity.Table) error {
	table.ID = uuid.New().String()
	table.CreatedAt = time.Now()
	table.UpdatedAt = time.Now()
	if table.Status == "" {
		table.Status = entity.TableStatusAvailable
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tables (id, name, capacity, status, created_at, updated_at) 
         VALUES (?, ?, ?, ?, ?, ?)`,
		table.ID, table.Name, table.Capacity, table.Status, table.CreatedAt, table.UpdatedAt)
	return err
}

func (r *tableRepository) Update(ctx context.Context, table *entity.Table) error {
	table.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`UPDATE tables SET name = ?, capacity = ?, updated_at = ? 
         WHERE id = ? AND deleted_at IS NULL`,
		table.Name, table.Capacity, table.UpdatedAt, table.ID)
	return err
}

func (r *tableRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE tables SET deleted_at = NOW() 
         WHERE id = ? AND deleted_at IS NULL`, id)
	return err
}

func (r *tableRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE tables SET status = ?, updated_at = ? 
         WHERE id = ? AND deleted_at IS NULL`, status, now, id)
	return err
}
