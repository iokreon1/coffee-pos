package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Category, error)
	FindByName(ctx context.Context, name string) (*entity.Category, error)
	FindAll(ctx context.Context) ([]entity.Category, error)
	Create(ctx context.Context, category *entity.Category) error
	Update(ctx context.Context, category *entity.Category) error
	Delete(ctx context.Context, id string) error
}

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*entity.Category, error) {
	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at, deleted_at 
         FROM categories WHERE id = ? AND deleted_at IS NULL`, id).
		Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) FindByName(ctx context.Context, name string) (*entity.Category, error) {
	category := &entity.Category{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, created_at, updated_at, deleted_at 
         FROM categories WHERE name = ? AND deleted_at IS NULL`, name).
		Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]entity.Category, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, created_at, updated_at, deleted_at 
         FROM categories WHERE deleted_at IS NULL ORDER BY name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []entity.Category{}
	for rows.Next() {
		var category entity.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt, &category.DeletedAt); err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, rows.Err()
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) error {
	category.ID = uuid.New().String()
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO categories (id, name, created_at, updated_at) 
         VALUES (?, ?, ?, ?)`,
		category.ID, category.Name, category.CreatedAt, category.UpdatedAt)
	return err
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) error {
	category.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET name = ?, updated_at = ? 
         WHERE id = ? AND deleted_at IS NULL`,
		category.Name, category.UpdatedAt, category.ID)
	return err
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE categories SET deleted_at = NOW() 
         WHERE id = ? AND deleted_at IS NULL`, id)
	return err
}
