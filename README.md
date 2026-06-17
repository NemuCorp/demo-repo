# Demo Store

A full-stack ecommerce application built with React + TypeScript on the frontend and Go + Gin on the backend, backed by PostgreSQL. Browse products, manage a shopping cart, and administer your product catalog — all through a clean, responsive web interface.

---

## Screenshots

### Home Page — Product Catalog
![Home Page](screenshots/home-page.svg)

Browse all available products in a responsive grid layout with pricing and stock information.

### Product Detail Page
![Product Detail](screenshots/product-page.svg)

View detailed product information, including description, price, and stock availability. Add items to your cart directly from this page.

### Shopping Cart
![Cart](screenshots/cart-page.svg)

Manage your cart with per-item quantity controls, subtotals, and a running total. Requires authentication.

### Login Page
![Login](screenshots/login-page.svg)

Sign in to your account to access cart functionality and the admin dashboard.

### Register Page
![Register](screenshots/register-page.svg)

Create a new account with email and password (minimum 6 characters).

### Admin Dashboard
![Admin Dashboard](screenshots/admin-dashboard.svg)

Overview of your store with quick access to product management.

### Admin — Product Management
![Admin Products](screenshots/admin-products.svg)

View all products in a table and create new products via a form interface.

---

## Features

- **Product Catalog** — Browse all products with a responsive grid layout
- **Product Details** — View detailed information for each product (name, description, price, stock)
- **User Authentication** — Register, login, and logout with bcrypt-hashed passwords and session-based Bearer token auth
- **Shopping Cart** — Add, update, and remove items; view per-item subtotals and cart total
- **Admin Panel** — Dashboard with product management: view products in a table, create new products
- **CLI Database Management** — Run migrations, rollbacks, and clean the database from the command line

---

## Use Cases

| Use Case | Description |
|----------|-------------|
| **Online Storefront** | List and sell products with a public-facing catalog |
| **Inventory Management** | Track product stock levels and manage listings via the admin panel |
| **Proof of Concept** | Use as a starting template for full-stack Go + React projects |
| **Learning Resource** | Demonstrates session-based auth, raw SQL with PostgreSQL, and React + TypeScript best practices |

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 18, TypeScript, React Router v6 |
| Backend | Go 1.21, Gin framework |
| Database | PostgreSQL (raw SQL with prepared statements) |
| Auth | bcrypt + SHA-256 session tokens |
| Styling | Plain CSS (responsive) |

---

## Getting Started

### Prerequisites

- **Go** 1.21 or newer
- **Node.js** 18 or newer with npm
- **PostgreSQL** running locally

### 1. Clone the Repository

```bash
git clone https://github.com/NemuCorp/demo-repo.git
cd demo-repo
```

### 2. Set Up the Database

Create the database and run the initial migration:

```bash
# Create the database
createdb demorepo

# Apply the migration
psql "postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable" \
  < server/db/migrations/001_initial.sql
```

Alternatively, start the backend server first and use the CLI:

```bash
cd server
go run . up     # Run all pending migrations
go run . down   # Rollback the last migration
go run . clean  # Drop all tables
```

### 3. Start the Backend

```bash
cd server

# Optional: set environment variables (defaults shown)
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable"
export PORT="8080"

# Start the server
go run .
```

The API server starts on `http://localhost:8080`.

### 4. Start the Frontend

```bash
cd client

# Install dependencies
npm install

# Start the development server
npm start
```

The React dev server starts on `http://localhost:3000` and proxies API requests to the backend.

### 5. Open the App

Navigate to [http://localhost:3000](http://localhost:3000) to start browsing.

---

## Project Structure

```
demo-repo/
├── client/                   React + TypeScript frontend
│   ├── src/
│   │   ├── components/       Reusable UI: Navbar, ProductCard, ProtectedRoute
│   │   ├── contexts/         AuthContext (authentication state management)
│   │   ├── pages/            Route-level page components
│   │   │   └── admin/        Admin dashboard and product management
│   │   └── services/         API client functions
│   └── package.json          Dependencies and proxy configuration
├── server/                   Go + Gin backend
│   ├── handler/              HTTP request handlers (auth, product, cart)
│   ├── db/                   Database access layer and migrations
│   ├── myerrors/             Sentinel error types
│   └── logger/               Logging utilities
└── screenshots/              Page screenshots
```

---

## Available Scripts

### Client (`client/`)

| Command | Description |
|---------|-------------|
| `npm start` | Start the React development server |
| `npm run build` | Create an optimized production build |
| `npm test` | Run the test suite |

### Server (`server/`)

| Command | Description |
|---------|-------------|
| `go run .` | Start the API server |
| `go run . up` | Run database migrations |
| `go run . down` | Rollback the last migration |
| `go run . clean` | Drop all database tables |
| `go test ./...` | Run all Go tests |
