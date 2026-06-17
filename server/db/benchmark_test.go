package db

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func openTestDB(tb testing.TB) *sql.DB {
	tb.Helper()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable"
	}
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		tb.Skipf("skipping benchmark: unable to open db: %v", err)
	}
	if err := conn.Ping(); err != nil {
		tb.Skipf("skipping benchmark: unable to ping db: %v", err)
	}
	return conn
}

func seedDB(tb testing.TB, conn *sql.DB) (userID int, productID int, cleanup func()) {
	tb.Helper()
	_, err := conn.Exec(`DELETE FROM cart_items`)
	if err != nil {
		tb.Fatalf("seed cleanup cart_items: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM sessions`)
	if err != nil {
		tb.Fatalf("seed cleanup sessions: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM products`)
	if err != nil {
		tb.Fatalf("seed cleanup products: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM users`)
	if err != nil {
		tb.Fatalf("seed cleanup users: %v", err)
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.MinCost)
	err = conn.QueryRow(
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		fmt.Sprintf("bench-%d@test.com", time.Now().UnixNano()), string(hash),
	).Scan(&userID)
	if err != nil {
		tb.Fatalf("seed user: %v", err)
	}

	err = conn.QueryRow(
		`INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id`,
		"Bench Product", "A benchmark product", 9.99, 100,
	).Scan(&productID)
	if err != nil {
		tb.Fatalf("seed product: %v", err)
	}

	cleanup = func() {
		conn.Exec(`DELETE FROM cart_items`)
		conn.Exec(`DELETE FROM sessions`)
		conn.Exec(`DELETE FROM products`)
		conn.Exec(`DELETE FROM users`)
	}

	return userID, productID, cleanup
}

func BenchmarkAuthGetUserByID(b *testing.B) {
	conn := openTestDB(b)
	auth, err := NewAuthDB(conn)
	if err != nil {
		b.Fatalf("NewAuthDB: %v", err)
	}
	userID, _, cleanup := seedDB(b, conn)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.GetUserByID(userID)
		if err != nil {
			b.Fatalf("GetUserByID: %v", err)
		}
	}
}

func BenchmarkAuthGetSession(b *testing.B) {
	conn := openTestDB(b)
	auth, err := NewAuthDB(conn)
	if err != nil {
		b.Fatalf("NewAuthDB: %v", err)
	}
	userID, _, cleanup := seedDB(b, conn)
	defer cleanup()

	sessionHash := uuid.New().String()
	_, err = auth.CreateSession(userID, sessionHash, time.Now().Add(1*time.Hour))
	if err != nil {
		b.Fatalf("CreateSession: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := auth.GetSession(sessionHash)
		if err != nil {
			b.Fatalf("GetSession: %v", err)
		}
	}
}

func BenchmarkAuthDeleteSession(b *testing.B) {
	conn := openTestDB(b)
	auth, err := NewAuthDB(conn)
	if err != nil {
		b.Fatalf("NewAuthDB: %v", err)
	}
	userID, _, cleanup := seedDB(b, conn)
	defer cleanup()

	sessionHashes := make([]uuid.UUID, b.N)
	for i := 0; i < b.N; i++ {
		s, err := auth.CreateSession(userID, uuid.New().String(), time.Now().Add(1*time.Hour))
		if err != nil {
			b.Fatalf("CreateSession: %v", err)
		}
		sessionHashes[i] = s.ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := auth.DeleteSession(sessionHashes[i])
		if err != nil {
			b.Fatalf("DeleteSession: %v", err)
		}
	}
}

func BenchmarkProductGetProductByID(b *testing.B) {
	conn := openTestDB(b)
	productDB, err := NewProductDB(conn)
	if err != nil {
		b.Fatalf("NewProductDB: %v", err)
	}
	_, productID, cleanup := seedDB(b, conn)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := productDB.GetProductByID(productID)
		if err != nil {
			b.Fatalf("GetProductByID: %v", err)
		}
	}
}

func BenchmarkProductListProducts(b *testing.B) {
	conn := openTestDB(b)
	productDB, err := NewProductDB(conn)
	if err != nil {
		b.Fatalf("NewProductDB: %v", err)
	}
	_, _, cleanup := seedDB(b, conn)
	defer cleanup()

	for i := 0; i < 100; i++ {
		_, err := productDB.CreateProduct(
			fmt.Sprintf("Product %d", i),
			fmt.Sprintf("Description %d", i),
			float64(i)+0.99,
			"/img/default.png",
			50,
		)
		if err != nil {
			b.Fatalf("CreateProduct: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := productDB.ListProducts()
		if err != nil {
			b.Fatalf("ListProducts: %v", err)
		}
	}
}

func BenchmarkProductListProductsPaginated(b *testing.B) {
	conn := openTestDB(b)
	productDB, err := NewProductDB(conn)
	if err != nil {
		b.Fatalf("NewProductDB: %v", err)
	}
	_, _, cleanup := seedDB(b, conn)
	defer cleanup()

	for i := 0; i < 100; i++ {
		_, err := productDB.CreateProduct(
			fmt.Sprintf("Product %d", i),
			fmt.Sprintf("Description %d", i),
			float64(i)+0.99,
			"/img/default.png",
			50,
		)
		if err != nil {
			b.Fatalf("CreateProduct: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := productDB.ListProductsPaginated(20, 0)
		if err != nil {
			b.Fatalf("ListProductsPaginated: %v", err)
		}
	}
}

func BenchmarkCartGetCart(b *testing.B) {
	conn := openTestDB(b)
	cartDB, err := NewCartDB(conn)
	if err != nil {
		b.Fatalf("NewCartDB: %v", err)
	}
	userID, productID, cleanup := seedDB(b, conn)
	defer cleanup()

	_, err = cartDB.AddItem(userID, productID, 1)
	if err != nil {
		b.Fatalf("AddItem: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cartDB.GetCart(userID)
		if err != nil {
			b.Fatalf("GetCart: %v", err)
		}
	}
}

func BenchmarkCartAddItem(b *testing.B) {
	conn := openTestDB(b)
	cartDB, err := NewCartDB(conn)
	if err != nil {
		b.Fatalf("NewCartDB: %v", err)
	}
	userID, productID, cleanup := seedDB(b, conn)
	defer cleanup()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := cartDB.AddItem(userID, productID, 1)
		if err != nil {
			b.Fatalf("AddItem: %v", err)
		}
	}
}

func BenchmarkCartRemoveItem(b *testing.B) {
	conn := openTestDB(b)
	cartDB, err := NewCartDB(conn)
	if err != nil {
		b.Fatalf("NewCartDB: %v", err)
	}
	userID, productID, cleanup := seedDB(b, conn)
	defer cleanup()

	for i := 0; i < b.N; i++ {
		_, err := cartDB.AddItem(userID, productID, 1)
		if err != nil {
			b.Fatalf("AddItem: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := cartDB.RemoveItem(userID, productID)
		if err != nil {
			b.Fatalf("RemoveItem: %v", err)
		}
	}
}
