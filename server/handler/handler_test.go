package handler

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/NemuCorp/demo-repo/server/db"
	"github.com/NemuCorp/demo-repo/server/myerrors"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable"
	}
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Skipf("skipping test: unable to open db: %v", err)
	}
	if err := conn.Ping(); err != nil {
		t.Skipf("skipping test: unable to ping db: %v", err)
	}
	return conn
}

func setupTestDB(t *testing.T) (*db.DB, func()) {
	t.Helper()
	conn := openTestDB(t)

	_, err := conn.Exec(`DELETE FROM cart_items`)
	if err != nil {
		t.Fatalf("cleanup cart_items: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM sessions`)
	if err != nil {
		t.Fatalf("cleanup sessions: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM products`)
	if err != nil {
		t.Fatalf("cleanup products: %v", err)
	}
	_, err = conn.Exec(`DELETE FROM users`)
	if err != nil {
		t.Fatalf("cleanup users: %v", err)
	}

	database := &db.DB{Conn: conn}
	if err := database.PrepareStatements(); err != nil {
		t.Fatalf("prepare statements: %v", err)
	}

	cleanup := func() {
		database.Close()
	}

	return database, cleanup
}

func createUser(t *testing.T, authDB *db.AuthDB, email, password string, isAdmin bool) *db.User {
	t.Helper()
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	user, err := authDB.CreateUser(email, string(hash), isAdmin)
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	return user
}

func createSession(t *testing.T, authDB *db.AuthDB, userID int) (plainToken string) {
	t.Helper()
	token := uuid.New().String()
	hash := sha256Sum(token)
	_, err := authDB.CreateSession(userID, hash, time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	return token
}

func sha256Sum(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func newTestEngine() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestAdminMiddleware(t *testing.T) {
	tests := []struct {
		name       string
		setAdmin   bool
		adminValue interface{}
		wantStatus int
	}{
		{
			name:       "allows admin user",
			setAdmin:   true,
			adminValue: true,
			wantStatus: http.StatusOK,
		},
		{
			name:       "forbids non-admin user",
			setAdmin:   true,
			adminValue: false,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "forbids missing is_admin flag",
			setAdmin:   false,
			adminValue: nil,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "forbids invalid is_admin type",
			setAdmin:   true,
			adminValue: "yes",
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestEngine()
			r.Use(func(c *gin.Context) {
				if tt.setAdmin {
					c.Set("is_admin", tt.adminValue)
				}
				c.Next()
			})
			r.GET("/admin", AdminMiddleware(), func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest(http.MethodGet, "/admin", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("AdminMiddleware() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestAuthHandler_Me(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	handler := NewAuthHandler(database.Auth)
	user := createUser(t, database.Auth, "me@test.com", "password123", false)
	token := createSession(t, database.Auth, user.ID)

	tests := []struct {
		name       string
		token      string
		wantStatus int
		wantAdmin  bool
	}{
		{
			name:       "returns authenticated user",
			token:      token,
			wantStatus: http.StatusOK,
			wantAdmin:  false,
		},
		{
			name:       "rejects missing token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestEngine()
			authMiddleware := AuthMiddleware(database.Auth)
			r.GET("/api/auth/me", authMiddleware, handler.Me)

			req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Me() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.wantStatus == http.StatusOK {
				var body struct {
					User struct {
						IsAdmin bool `json:"is_admin"`
					} `json:"user"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if body.User.IsAdmin != tt.wantAdmin {
					t.Errorf("Me() is_admin = %v, want %v", body.User.IsAdmin, tt.wantAdmin)
				}
			}
		})
	}
}

func TestAuthHandler_Me_Admin(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	handler := NewAuthHandler(database.Auth)
	user := createUser(t, database.Auth, "admin-me@test.com", "password123", true)
	token := createSession(t, database.Auth, user.ID)

	r := newTestEngine()
	authMiddleware := AuthMiddleware(database.Auth)
	r.GET("/api/auth/me", authMiddleware, handler.Me)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Me() status = %d, want %d", w.Code, http.StatusOK)
	}

	var body struct {
		User struct {
			IsAdmin bool `json:"is_admin"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if !body.User.IsAdmin {
		t.Errorf("Me() is_admin = false, want true")
	}
}

func TestProductHandler_Update(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	handler := NewProductHandler(database.Product)
	admin := createUser(t, database.Auth, "admin-product@test.com", "password123", true)
	adminToken := createSession(t, database.Auth, admin.ID)

	product, err := database.Product.CreateProduct("Original", "Original desc", 9.99, "/original.png", 10)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	tests := []struct {
		name       string
		id         string
		body       map[string]interface{}
		token      string
		wantStatus int
	}{
		{
			name:       "updates existing product as admin",
			id:         fmt.Sprintf("%d", product.ID),
			body:       map[string]interface{}{"name": "Updated", "description": "Updated desc", "price": 19.99, "image_path": "/updated.png", "stock": 20},
			token:      adminToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "returns not found for missing product",
			id:         "999999",
			body:       map[string]interface{}{"name": "Updated", "description": "Updated desc", "price": 19.99, "image_path": "/updated.png", "stock": 20},
			token:      adminToken,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "rejects invalid product id",
			id:         "abc",
			body:       map[string]interface{}{"name": "Updated", "description": "Updated desc", "price": 19.99, "image_path": "/updated.png", "stock": 20},
			token:      adminToken,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects missing auth",
			id:         fmt.Sprintf("%d", product.ID),
			body:       map[string]interface{}{"name": "Updated", "description": "Updated desc", "price": 19.99, "image_path": "/updated.png", "stock": 20},
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects non-admin user",
			id:         fmt.Sprintf("%d", product.ID),
			body:       map[string]interface{}{"name": "Updated", "description": "Updated desc", "price": 19.99, "image_path": "/updated.png", "stock": 20},
			token:      func() string { u := createUser(t, database.Auth, "user-product@test.com", "password123", false); return createSession(t, database.Auth, u.ID) }(),
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestEngine()
			authMiddleware := AuthMiddleware(database.Auth)
			adminMiddleware := AdminMiddleware()
			r.PUT("/api/products/:id", authMiddleware, adminMiddleware, handler.Update)

			bodyBytes, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/products/"+tt.id, bytes.NewReader(bodyBytes))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Update() status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestProductHandler_Delete(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	handler := NewProductHandler(database.Product)
	admin := createUser(t, database.Auth, "admin-delete@test.com", "password123", true)
	adminToken := createSession(t, database.Auth, admin.ID)

	product, err := database.Product.CreateProduct("To Delete", "Desc", 5.99, "/delete.png", 5)
	if err != nil {
		t.Fatalf("create product: %v", err)
	}

	tests := []struct {
		name       string
		id         string
		token      string
		wantStatus int
	}{
		{
			name:       "deletes existing product as admin",
			id:         fmt.Sprintf("%d", product.ID),
			token:      adminToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "returns not found when product does not exist",
			id:         "999999",
			token:      adminToken,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "rejects invalid product id",
			id:         "abc",
			token:      adminToken,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "rejects missing auth",
			id:         fmt.Sprintf("%d", product.ID),
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "rejects non-admin user",
			id:         fmt.Sprintf("%d", product.ID),
			token:      func() string { u := createUser(t, database.Auth, "user-delete@test.com", "password123", false); return createSession(t, database.Auth, u.ID) }(),
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := newTestEngine()
			authMiddleware := AuthMiddleware(database.Auth)
			adminMiddleware := AdminMiddleware()
			r.DELETE("/api/products/:id", authMiddleware, adminMiddleware, handler.Delete)

			req := httptest.NewRequest(http.MethodDelete, "/api/products/"+tt.id, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("Delete() status = %d, want %d", w.Code, tt.wantStatus)
			}

			if tt.name == "returns not found when product does not exist" {
				var body struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
					t.Fatalf("unmarshal response: %v", err)
				}
				if body.Error != myerrors.ErrProductNotFound.Error() {
					t.Errorf("Delete() error = %q, want %q", body.Error, myerrors.ErrProductNotFound.Error())
				}
			}
		})
	}
}
