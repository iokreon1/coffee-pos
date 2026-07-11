package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type ProductFilter struct {
	CategoryID string
	IsActive   *bool
	Search     string
	Page       int
	Limit      int
}

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Product, error)
	FindAll(ctx context.Context, filter ProductFilter) ([]entity.Product, int, error)
	Create(ctx context.Context, product *entity.Product) error
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id string) error
	UpdateStock(ctx context.Context, id string, stock int) error
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*entity.Product, error) {
	product := &entity.Product{}
	var catID sql.NullString
	var catName sql.NullString

	query := `
		SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock,
		       p.image_url, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		WHERE p.id = ? AND p.deleted_at IS NULL
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.CategoryID, &product.Name, &product.Description, &product.Price, &product.Stock,
		&product.ImageURL, &product.IsActive, &product.CreatedAt, &product.UpdatedAt,
		&catID, &catName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if catID.Valid {
		product.Category = &entity.Category{
			ID:   catID.String,
			Name: catName.String,
		}
	}

	return product, nil
}

func (r *productRepository) FindAll(ctx context.Context, filter ProductFilter) ([]entity.Product, int, error) {
	var whereClauses []string
	var args []interface{}

	whereClauses = append(whereClauses, "p.deleted_at IS NULL")

	if filter.CategoryID != "" {
		whereClauses = append(whereClauses, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}

	if filter.IsActive != nil {
		whereClauses = append(whereClauses, "p.is_active = ?")
		args = append(args, *filter.IsActive)
	}

	if filter.Search != "" {
		whereClauses = append(whereClauses, "p.name LIKE ?")
		args = append(args, "%"+filter.Search+"%")
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	// Count query
	countQuery := `SELECT COUNT(*) FROM products p ` + whereSQL
	var total int
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Limit and offset configuration
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
		SELECT p.id, p.category_id, p.name, p.description, p.price, p.stock,
		       p.image_url, p.is_active, p.created_at, p.updated_at,
		       c.id, c.name
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
	` + whereSQL + ` ORDER BY p.created_at DESC LIMIT ? OFFSET ?`

	selectArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := []entity.Product{}
	for rows.Next() {
		var p entity.Product
		var catID sql.NullString
		var catName sql.NullString

		err := rows.Scan(
			&p.ID, &p.CategoryID, &p.Name, &p.Description, &p.Price, &p.Stock,
			&p.ImageURL, &p.IsActive, &p.CreatedAt, &p.UpdatedAt,
			&catID, &catName,
		)
		if err != nil {
			return nil, 0, err
		}

		if catID.Valid {
			p.Category = &entity.Category{
				ID:   catID.String,
				Name: catName.String,
			}
		}

		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	product.ID = uuid.New().String()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	query := `
		INSERT INTO products (id, category_id, name, description, price, stock, image_url, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.CategoryID, product.Name, product.Description,
		product.Price, product.Stock, product.ImageURL, product.IsActive,
		product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	product.UpdatedAt = time.Now()

	query := `
		UPDATE products 
		SET name = ?, description = ?, price = ?, image_url = ?, is_active = ?, category_id = ?, updated_at = ?
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		product.Name, product.Description, product.Price, product.ImageURL,
		product.IsActive, product.CategoryID, product.UpdatedAt, product.ID,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE products 
		SET deleted_at = NOW() 
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepository) UpdateStock(ctx context.Context, id string, stock int) error {
	query := `
		UPDATE products 
		SET stock = ?, updated_at = ? 
		WHERE id = ? AND deleted_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, stock, time.Now(), id)
	return err
}
