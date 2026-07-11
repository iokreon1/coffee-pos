package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type UserRepository interface {
	FindByID(ctx context.Context, id string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context) ([]entity.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
         FROM users WHERE id = ? AND deleted_at IS NULL`, id).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
         FROM users WHERE email = ? AND deleted_at IS NULL`, email).
		Scan(&user.ID, &user.Name, &user.Email, &user.Password,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, password, role, is_active, created_at, updated_at)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.Password,
		user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		`UPDATE users SET name=?, email=?, role=?, is_active=?, updated_at=?
         WHERE id=? AND deleted_at IS NULL`,
		user.Name, user.Email, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	return err
}

func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, email, password, role, is_active, created_at, updated_at, deleted_at
         FROM users WHERE deleted_at IS NULL ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []entity.User{}
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}
