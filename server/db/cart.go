package db

import (
	"database/sql"
	"time"

	"github.com/NemuCorp/demo-repo/server/logger"
)

type CartDB struct {
	addItem      *sql.Stmt
	getCart      *sql.Stmt
	updateItem   *sql.Stmt
	removeItem   *sql.Stmt
	clearCart    *sql.Stmt
}

func NewCartDB(conn *sql.DB) (*CartDB, error) {
	var c CartDB

	stmt, err := conn.Prepare(`
		WITH inserted AS (
			INSERT INTO cart_items (user_id, product_id, quantity)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id, product_id)
			DO UPDATE SET quantity = cart_items.quantity + $3
			RETURNING id, user_id, product_id, quantity, created_at
		)
		SELECT i.id, i.user_id, i.product_id, i.quantity, i.created_at,
			   p.name, p.price, p.image_path
		FROM inserted i
		JOIN products p ON p.id = i.product_id
	`)
	if err != nil {
		return nil, err
	}
	c.addItem = stmt

	stmt, err = conn.Prepare(`
		SELECT ci.id, ci.user_id, ci.product_id, ci.quantity, ci.created_at,
			   p.name, p.price, p.image_path
		FROM cart_items ci
		JOIN products p ON p.id = ci.product_id
		WHERE ci.user_id = $1
	`)
	if err != nil {
		return nil, err
	}
	c.getCart = stmt

	stmt, err = conn.Prepare(`
		WITH updated AS (
			UPDATE cart_items SET quantity = $3 WHERE user_id = $1 AND product_id = $2
			RETURNING id, user_id, product_id, quantity, created_at
		)
		SELECT u.id, u.user_id, u.product_id, u.quantity, u.created_at,
			   p.name, p.price, p.image_path
		FROM updated u
		JOIN products p ON p.id = u.product_id
	`)
	if err != nil {
		return nil, err
	}
	c.updateItem = stmt

	stmt, err = conn.Prepare(`DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`)
	if err != nil {
		return nil, err
	}
	c.removeItem = stmt

	stmt, err = conn.Prepare(`DELETE FROM cart_items WHERE user_id = $1`)
	if err != nil {
		return nil, err
	}
	c.clearCart = stmt

	logger.Info.Println("CartDB prepared statements initialized")
	return &c, nil
}

type CartItem struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProductID   int       `json:"product_id"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	ProductName string    `json:"product_name"`
	Price       float64   `json:"price"`
	ImagePath   string    `json:"image_path,omitempty"`
}

func (c *CartDB) AddItem(userID, productID, quantity int) (*CartItem, error) {
	ci := &CartItem{}
	var imagePath sql.NullString
	err := c.addItem.QueryRow(userID, productID, quantity).Scan(
		&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt,
		&ci.ProductName, &ci.Price, &imagePath,
	)
	if err != nil {
		return nil, err
	}
	if imagePath.Valid {
		ci.ImagePath = imagePath.String
	}
	return ci, nil
}

func (c *CartDB) GetCart(userID int) ([]CartItem, error) {
	rows, err := c.getCart.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	for rows.Next() {
		var ci CartItem
		var imagePath sql.NullString
		if err := rows.Scan(&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt, &ci.ProductName, &ci.Price, &imagePath); err != nil {
			return nil, err
		}
		if imagePath.Valid {
			ci.ImagePath = imagePath.String
		}
		items = append(items, ci)
	}
	return items, rows.Err()
}

func (c *CartDB) UpdateItem(userID, productID, quantity int) (*CartItem, error) {
	ci := &CartItem{}
	var imagePath sql.NullString
	err := c.updateItem.QueryRow(userID, productID, quantity).Scan(
		&ci.ID, &ci.UserID, &ci.ProductID, &ci.Quantity, &ci.CreatedAt,
		&ci.ProductName, &ci.Price, &imagePath,
	)
	if err != nil {
		return nil, err
	}
	if imagePath.Valid {
		ci.ImagePath = imagePath.String
	}
	return ci, nil
}

func (c *CartDB) RemoveItem(userID, productID int) error {
	_, err := c.removeItem.Exec(userID, productID)
	return err
}

func (c *CartDB) ClearCart(userID int) error {
	_, err := c.clearCart.Exec(userID)
	return err
}
