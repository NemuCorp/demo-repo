package db

import (
	"database/sql"
	"fmt"
	"time"

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

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(5 * time.Minute)

	d := &DB{Conn: conn}

	authDB, err := NewAuthDB(conn)
	if err != nil {
		return nil, fmt.Errorf("auth db init: %w", err)
	}
	d.Auth = authDB

	cartDB, err := NewCartDB(conn)
	if err != nil {
		return nil, fmt.Errorf("cart db init: %w", err)
	}
	d.Cart = cartDB

	productDB, err := NewProductDB(conn)
	if err != nil {
		return nil, fmt.Errorf("product db init: %w", err)
	}
	d.Product = productDB

	return d, nil
}

func (d *DB) Close() error {
	return d.Conn.Close()
}
