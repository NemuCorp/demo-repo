package db

import (
	"database/sql"
	"time"

	"github.com/NemuCorp/demo-repo/server/logger"
)

type ProductDB struct {
	createProduct  *sql.Stmt
	getProductByID *sql.Stmt
	listProducts   *sql.Stmt
	updateProduct  *sql.Stmt
	deleteProduct  *sql.Stmt
}

func NewProductDB(conn *sql.DB) (*ProductDB, error) {
	var p ProductDB

	stmt, err := conn.Prepare(`
		INSERT INTO products (name, description, price, image_path, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, price, image_path, stock, created_at, updated_at
	`)
	if err != nil {
		return nil, err
	}
	p.createProduct = stmt

	stmt, err = conn.Prepare(`
		SELECT id, name, description, price, image_path, stock, created_at, updated_at
		FROM products WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	p.getProductByID = stmt

	stmt, err = conn.Prepare(`
		SELECT id, name, description, price, image_path, stock, created_at, updated_at
		FROM products ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	p.listProducts = stmt

	stmt, err = conn.Prepare(`
		UPDATE products SET name = $2, description = $3, price = $4, image_path = $5, stock = $6, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, description, price, image_path, stock, created_at, updated_at
	`)
	if err != nil {
		return nil, err
	}
	p.updateProduct = stmt

	stmt, err = conn.Prepare(`DELETE FROM products WHERE id = $1`)
	if err != nil {
		return nil, err
	}
	p.deleteProduct = stmt

	logger.Info.Println("ProductDB prepared statements initialized")
	return &p, nil
}

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Price       float64   `json:"price"`
	ImagePath   string    `json:"image_path,omitempty"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *ProductDB) CreateProduct(name, description string, price float64, imagePath string, stock int) (*Product, error) {
	prod := &Product{}
	err := p.createProduct.QueryRow(name, description, price, imagePath, stock).Scan(
		&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.ImagePath, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (p *ProductDB) GetProductByID(id int) (*Product, error) {
	prod := &Product{}
	var desc, imagePath sql.NullString
	err := p.getProductByID.QueryRow(id).Scan(
		&prod.ID, &prod.Name, &desc, &prod.Price, &imagePath, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if desc.Valid {
		prod.Description = desc.String
	}
	if imagePath.Valid {
		prod.ImagePath = imagePath.String
	}
	return prod, nil
}

func (p *ProductDB) ListProducts() ([]Product, error) {
	rows, err := p.listProducts.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var prod Product
		var desc, imagePath sql.NullString
		if err := rows.Scan(&prod.ID, &prod.Name, &desc, &prod.Price, &imagePath, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt); err != nil {
			return nil, err
		}
		if desc.Valid {
			prod.Description = desc.String
		}
		if imagePath.Valid {
			prod.ImagePath = imagePath.String
		}
		products = append(products, prod)
	}
	return products, rows.Err()
}

func (p *ProductDB) UpdateProduct(id int, name, description string, price float64, imagePath string, stock int) (*Product, error) {
	prod := &Product{}
	err := p.updateProduct.QueryRow(id, name, description, price, imagePath, stock).Scan(
		&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.ImagePath, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

func (p *ProductDB) DeleteProduct(id int) error {
	_, err := p.deleteProduct.Exec(id)
	return err
}
