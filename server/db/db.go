package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn   *sql.DB
	Auth   *AuthDB
	Cart   *CartDB
	Product *ProductDB
}

func Init(connStr string) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	db := &DB{Conn: conn}

	authDB, err := NewAuthDB(conn)
	if err != nil {
		return nil, fmt.Errorf("auth db init: %w", err)
	}
	db.Auth = authDB

	cartDB, err := NewCartDB(conn)
	if err != nil {
		return nil, fmt.Errorf("cart db init: %w", err)
	}
	db.Cart = cartDB

	productDB, err := NewProductDB(conn)
	if err != nil {
		return nil, fmt.Errorf("product db init: %w", err)
	}
	db.Product = productDB

	return db, nil
}

func (d *DB) Close() error {
	return d.Conn.Close()
}
