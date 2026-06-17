package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DB struct {
	Conn     *sql.DB
	Auth     *AuthDB
	Cart     *CartDB
	Product  *ProductDB
	Tracking *TrackingDB
}

func Open(connStr string) (*DB, error) {
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("db open: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("db ping: %w", err)
	}

	return &DB{Conn: conn}, nil
}

func (db *DB) PrepareStatements() error {
	authDB, err := NewAuthDB(db.Conn)
	if err != nil {
		return fmt.Errorf("auth db init: %w", err)
	}
	db.Auth = authDB

	cartDB, err := NewCartDB(db.Conn)
	if err != nil {
		return fmt.Errorf("cart db init: %w", err)
	}
	db.Cart = cartDB

	productDB, err := NewProductDB(db.Conn)
	if err != nil {
		return fmt.Errorf("product db init: %w", err)
	}
	db.Product = productDB

	trackingDB, err := NewTrackingDB(db.Conn)
	if err != nil {
		return fmt.Errorf("tracking db init: %w", err)
	}
	db.Tracking = trackingDB

	return nil
}

func Init(connStr string) (*DB, error) {
	db, err := Open(connStr)
	if err != nil {
		return nil, err
	}

	if err := db.PrepareStatements(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func (d *DB) Close() error {
	return d.Conn.Close()
}
