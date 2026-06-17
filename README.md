# DemoRepo

A full-stack e-commerce web application built with React, TypeScript, Go, and PostgreSQL. Browse products, manage a shopping cart, and handle user authentication — all in one clean monorepo.

---

## Screenshots

### Home Page — Product Catalog

![Home Page](screenshots/home-page.svg)

Browse all available products in a responsive grid layout. Each product card shows the name, price, stock status, and an image placeholder.

### Product Detail Page

![Product Detail](screenshots/product-page.svg)

View detailed information about a product, including its description, price, stock availability, and an image. Authenticated users can add items to their cart directly from this page.

### Shopping Cart

![Cart](screenshots/cart-page.svg)

Review and manage items in your cart. Update quantities or remove items before proceeding to checkout.

### Login & Register

| Login | Register |
|-------|----------|
| ![Login](screenshots/login-page.svg) | ![Register](screenshots/register-page.svg) |

Secure authentication with hashed passwords and session tokens. Users can register a new account or log in with existing credentials.

### Admin Dashboard

![Admin Dashboard](screenshots/admin-dashboard.svg)

Administrators can manage products — create, edit, and delete listings — from a dedicated admin panel.

## Features

### For Customers
- **Product Catalog** — Browse all available products with name, price, stock, and images.
- **Product Details** — View in-depth product information, including descriptions and stock levels.
- **Shopping Cart** — Add, update, and remove items from a persistent cart tied to your account.
- **User Authentication** — Register an account, log in, and log out securely with bcrypt-hashed passwords.
- **Multiple Sessions** — Stay logged in across multiple devices with session token support.

### For Administrators
- **Admin Dashboard** — Protected admin area for managing store content.
- **Product Management** — Create new products, edit existing ones, and remove discontinued items.

### Use Cases
- **Small to medium online stores** — Launch a product catalog with shopping cart functionality quickly.
- **Learning full-stack development** — Study a real-world monorepo combining React, Go, and PostgreSQL.
- **Bootstrapping an e-commerce MVP** — Use the project as a starting point for a custom online shop.
- **Testing frontend/backend integration** — Explore how a React SPA communicates with a Go REST API.

## Tech Stack

| Layer    | Technology                           |
|----------|--------------------------------------|
| Frontend | React 18, TypeScript, React Router 6 |
| Backend  | Go 1.21, Gin web framework           |
| Database | PostgreSQL, raw SQL (no ORM)         |
| Auth     | bcrypt password hashing, session tokens |
| Logging  | `log/slog` (structured logging)      |

## Project Structure

```
demo-repo/
├── client/                 # React + TypeScript frontend
│   ├── public/             # Static assets
│   └── src/
│       ├── components/     # Reusable UI components (Navbar, ProductCard, ProtectedRoute)
│       ├── contexts/       # React contexts (AuthContext)
│       ├── pages/          # Page components (Home, Product, Cart, Login, Register, Admin)
│       ├── services/       # API client and data-fetching logic
│       ├── App.tsx         # Root component with routes
│       └── index.tsx       # Entry point
├── server/                 # Go + Gin backend
│   ├── db/                 # Database access, prepared statements, migrations
│   ├── handler/            # HTTP handlers (auth, cart, product)
│   ├── logger/             # Structured logging setup
│   ├── myerrors/           # Sentinel error types
│   ├── main.go             # Server entry point and route configuration
│   └── cmd.go              # CLI commands (migrate, clean, import, export)
└── screenshots/            # Application screenshots
```

## Setup & Running

### Quick Start with Docker Compose

The fastest way to run the entire stack is with a single command:

```bash
git clone https://github.com/NemuCorp/demo-repo.git
cd demo-repo
docker compose up
```

This builds and starts all three services (PostgreSQL, Go backend, React frontend). Once running:

- Frontend: [http://localhost:3000](http://localhost:3000)
- API: [http://localhost:8080](http://localhost:8080)

Database migrations run automatically on startup. To stop: `docker compose down`.

### Manual Setup

#### Prerequisites

- **Go 1.21+** — [Download](https://go.dev/dl/)
- **Node.js 18+** and **npm** — [Download](https://nodejs.org/)
- **PostgreSQL** — Running instance (local or remote)

### 1. Clone the repository

```bash
git clone https://github.com/NemuCorp/demo-repo.git
cd demo-repo
```

### 2. Set up the database

Make sure PostgreSQL is running, then set the connection string:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/demorepo?sslmode=disable"
```

Run database migrations (from the `server/` directory):

```bash
cd server
go run . up
```

> **Available CLI commands:** `up` (run migrations), `down` (rollback), `clean` (drop all tables), `import`, `export`

### 3. Start the backend server

```bash
cd server
go run .
```

The API server starts on `http://localhost:8080`.

### 4. Start the frontend client

In a separate terminal:

```bash
cd client
npm install
npm start
```

The React dev server starts on `http://localhost:3000` and proxies API requests to the Go backend.

### Environment Variables

| Variable       | Default               | Description                  |
|----------------|-----------------------|------------------------------|
| `DATABASE_URL` | (see above)           | PostgreSQL connection string |
| `PORT`         | `8080`                | Backend server port          |

## API Endpoints

### Authentication
| Method | Path             | Auth     | Description        |
|--------|------------------|----------|--------------------|
| POST   | `/api/auth/register` | No   | Create a new account |
| POST   | `/api/auth/login`    | No   | Log in, get session   |
| POST   | `/api/auth/logout`   | Yes  | End current session   |

### Products
| Method | Path              | Auth | Description            |
|--------|-------------------|------|------------------------|
| GET    | `/api/products`   | No   | List all products      |
| GET    | `/api/products/:id` | No | Get product by ID      |
| POST   | `/api/products`   | No   | Create a new product   |

### Cart
| Method | Path                    | Auth | Description              |
|--------|-------------------------|------|--------------------------|
| GET    | `/api/cart`             | Yes  | View current user's cart |
| POST   | `/api/cart`             | Yes  | Add item to cart         |
| PUT    | `/api/cart/:productId`  | Yes  | Update item quantity     |
| DELETE | `/api/cart/:productId`  | Yes  | Remove item from cart    |

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
