package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iokreon1/coffee-pos/internal/entity"
)

type OrderFilter struct {
	ShiftID   string
	CashierID string
	Status    string
	Page      int
	Limit     int
}

type OrderRepository interface {
	FindByID(ctx context.Context, id string) (*entity.Order, error)
	FindAll(ctx context.Context, filter OrderFilter) ([]entity.Order, int, error)
	Create(ctx context.Context, order *entity.Order) error
	Update(ctx context.Context, order *entity.Order) error
	AddItem(ctx context.Context, item *entity.OrderItem) error
	UpdateItem(ctx context.Context, item *entity.OrderItem) error
	DeleteItem(ctx context.Context, itemID string) error
	FindItemByID(ctx context.Context, itemID string) (*entity.OrderItem, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByID(ctx context.Context, id string) (*entity.Order, error) {
	order := &entity.Order{}
	var cashierID string
	var cashierName string
	var tableID sql.NullString
	var tableName sql.NullString

	query := `
		SELECT o.id, o.shift_id, o.cashier_id, o.table_id, o.promo_id, o.status,
		       o.subtotal, o.discount_amount, o.total, o.notes, o.created_at, o.updated_at,
		       u.id, u.name,
		       t.id, t.name
		FROM orders o
		JOIN users u ON u.id = o.cashier_id
		LEFT JOIN tables t ON t.id = o.table_id AND t.deleted_at IS NULL
		WHERE o.id = ?
	`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID, &order.ShiftID, &order.CashierID, &order.TableID, &order.PromoID, &order.Status,
		&order.Subtotal, &order.DiscountAmount, &order.Total, &order.Notes, &order.CreatedAt, &order.UpdatedAt,
		&cashierID, &cashierName,
		&tableID, &tableName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	order.Cashier = &entity.User{
		ID:   cashierID,
		Name: cashierName,
	}

	if tableID.Valid {
		order.Table = &entity.Table{
			ID:   tableID.String,
			Name: tableName.String,
		}
	}

	// Fetch items
	itemsQuery := `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.subtotal,
		       oi.notes, oi.created_at, oi.updated_at,
		       p.id, p.name, p.price
		FROM order_items oi
		JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = ?
	`

	rows, err := r.db.QueryContext(ctx, itemsQuery, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []entity.OrderItem{}
	for rows.Next() {
		var item entity.OrderItem
		var p entity.Product
		err := rows.Scan(
			&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.Subtotal,
			&item.Notes, &item.CreatedAt, &item.UpdatedAt,
			&p.ID, &p.Name, &p.Price,
		)
		if err != nil {
			return nil, err
		}
		item.Product = &p
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	order.Items = items
	return order, nil
}

func (r *orderRepository) FindAll(ctx context.Context, filter OrderFilter) ([]entity.Order, int, error) {
	var whereClauses []string
	var args []interface{}

	if filter.ShiftID != "" {
		whereClauses = append(whereClauses, "o.shift_id = ?")
		args = append(args, filter.ShiftID)
	}
	if filter.CashierID != "" {
		whereClauses = append(whereClauses, "o.cashier_id = ?")
		args = append(args, filter.CashierID)
	}
	if filter.Status != "" {
		whereClauses = append(whereClauses, "o.status = ?")
		args = append(args, filter.Status)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total orders
	countQuery := `SELECT COUNT(*) FROM orders o ` + whereSQL
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
		SELECT o.id, o.shift_id, o.cashier_id, o.table_id, o.promo_id, o.status,
		       o.subtotal, o.discount_amount, o.total, o.notes, o.created_at, o.updated_at,
		       u.id, u.name
		FROM orders o
		JOIN users u ON u.id = o.cashier_id
	` + whereSQL + ` ORDER BY o.created_at DESC LIMIT ? OFFSET ?`

	selectArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, selectSQL, selectArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := []entity.Order{}
	for rows.Next() {
		var o entity.Order
		var u entity.User
		err := rows.Scan(
			&o.ID, &o.ShiftID, &o.CashierID, &o.TableID, &o.PromoID, &o.Status,
			&o.Subtotal, &o.DiscountAmount, &o.Total, &o.Notes, &o.CreatedAt, &o.UpdatedAt,
			&u.ID, &u.Name,
		)
		if err != nil {
			return nil, 0, err
		}
		o.Cashier = &u
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) Create(ctx context.Context, order *entity.Order) error {
	order.ID = uuid.New().String()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	query := `
		INSERT INTO orders (
			id, shift_id, cashier_id, table_id, promo_id, status,
			subtotal, discount_amount, total, notes, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		order.ID, order.ShiftID, order.CashierID, order.TableID, order.PromoID, order.Status,
		order.Subtotal, order.DiscountAmount, order.Total, order.Notes, order.CreatedAt, order.UpdatedAt,
	)
	return err
}

func (r *orderRepository) Update(ctx context.Context, order *entity.Order) error {
	order.UpdatedAt = time.Now()

	query := `
		UPDATE orders 
		SET table_id = ?, promo_id = ?, status = ?, subtotal = ?, discount_amount = ?, total = ?, notes = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		order.TableID, order.PromoID, order.Status, order.Subtotal, order.DiscountAmount, order.Total, order.Notes, order.UpdatedAt,
		order.ID,
	)
	return err
}

func (r *orderRepository) AddItem(ctx context.Context, item *entity.OrderItem) error {
	item.ID = uuid.New().String()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	item.Subtotal = item.Price * int64(item.Quantity)

	query := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price, subtotal, notes, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.OrderID, item.ProductID, item.Quantity, item.Price, item.Subtotal, item.Notes, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *orderRepository) UpdateItem(ctx context.Context, item *entity.OrderItem) error {
	item.UpdatedAt = time.Now()

	query := `
		UPDATE order_items 
		SET quantity = ?, notes = ?, subtotal = ?, updated_at = ? 
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		item.Quantity, item.Notes, item.Subtotal, item.UpdatedAt,
		item.ID,
	)
	return err
}

func (r *orderRepository) DeleteItem(ctx context.Context, itemID string) error {
	query := `DELETE FROM order_items WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}

func (r *orderRepository) FindItemByID(ctx context.Context, itemID string) (*entity.OrderItem, error) {
	item := &entity.OrderItem{}

	query := `
		SELECT id, order_id, product_id, quantity, price, subtotal, notes, created_at, updated_at
		FROM order_items 
		WHERE id = ?
	`

	err := r.db.QueryRowContext(ctx, query, itemID).Scan(
		&item.ID, &item.OrderID, &item.ProductID, &item.Quantity, &item.Price, &item.Subtotal, &item.Notes, &item.CreatedAt, &item.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return item, nil
}
