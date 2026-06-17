# Create a Suitable Frontend

**Issue:** [NemuCorp/demo-repo#7](https://github.com/NemuCorp/demo-repo/issues/7)

## Problem

The existing `client/` directory contains only a placeholder React app (`App.tsx` rendering "Client application is running") with no routing, no pages, and no integration with the backend API. The Go server already exposes a full ecommerce API (auth, products, cart) and the frontend needs to consume it. The frontend must serve two distinct audiences: public shoppers browsing products and managing their cart, and administrators managing products.

## Goals

- Add `react-router-dom` for client-side routing, splitting the app into public (client) and admin routes.
- Create reusable UI components (header/nav, product cards, cart item rows, forms) following existing convention.
- Build public-facing pages: Home, Product Listing, Product Detail, Cart, Login, Register.
- Build admin-facing pages: Dashboard and Product Management (list/create only).
- Create an API service layer (`src/services/`) to call the backend endpoints, with auth token management.
- Use existing `package.json` toolchain (react-scripts, TypeScript) without ejecting or changing the build system.
- Store the auth token in `localStorage` and attach it via an API helper; handle 401 responses by redirecting to login.
- All pages are functional against the live API endpoints documented in the server handlers.

## Approach

### Dependencies

Add `react-router-dom` (v6) to `client/package.json`. No other new runtime dependencies; use plain React state and `fetch` for API calls.

### Route Structure

```
/                    → Home page (product listing for shoppers)
/products/:id        → Product Detail page
/cart                → Cart page (auth required)
/login               → Login page
/register            → Register page
/admin               → Admin Dashboard
/admin/products      → Admin Product Management
```

A route guard component wraps protected routes (cart, admin) and redirects to `/login` if no auth token is present.

### Component Tree

```
src/
  components/
    Header.tsx          # Navigation bar (links change based on auth state / admin)
    ProductCard.tsx      # Product thumbnail card used on Home and Product Detail
    CartItemRow.tsx      # Single cart item row with quantity controls and remove
    ProtectedRoute.tsx   # Auth guard wrapper that redirects to /login if unauthenticated
    AdminHeader.tsx      # Admin-specific navigation
  pages/
    Home.tsx             # Public product listing grid
    ProductDetail.tsx    # Single product view with add-to-cart
    Cart.tsx             # Cart view with update/remove actions
    Login.tsx            # Login form
    Register.tsx         # Registration form
    admin/
      Dashboard.tsx      # Admin landing with summary stats (placeholder)
      ProductList.tsx    # Admin product table (list + create only)
      ProductForm.tsx    # Create product form
  hooks/
    useAuth.tsx          # Auth context/hook for token state and user info
    useApi.tsx           # Generic fetch wrapper that injects auth header
  services/
    api.ts               # Base fetch helper (attach Bearer token, handle 401 → redirect)
    auth.ts              # login(), register(), logout() functions
    products.ts          # listProducts(), getProduct(), createProduct() functions
    cart.ts              # getCart(), addToCart(), updateCartItem(), removeCartItem()
  App.tsx                # Router setup with <BrowserRouter> and <Routes>
  index.tsx              # (unchanged, already renders <App />)
```

### API Integration

The Go server listens on port 8080 and exposes:

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | /api/auth/register | No | Create account |
| POST | /api/auth/login | No | Get session token |
| POST | /api/auth/logout | Yes | Invalidate sessions |
| GET | /api/products | No | List all products |
| GET | /api/products/:id | No | Get product by ID |
| POST | /api/products | Yes (admin) | Create product |
| GET | /api/cart | Yes | View cart |
| POST | /api/cart | Yes | Add item to cart |
| PUT | /api/cart/:productId | Yes | Update cart item quantity |
| DELETE | /api/cart/:productId | Yes | Remove item from cart |

#### Request/Response Schemas

**Auth endpoints:**

- `POST /api/auth/register` — request: `{ email, password }` (password min 6 chars); response: `{ user: { id, email } }`
- `POST /api/auth/login` — request: `{ email, password }`; response: `{ token, session: { id, user_id, created_at, expires_at } }`
- `POST /api/auth/logout` — request: (none, token in header); response: `{ message }`

**Product endpoints:**

- `GET /api/products` — response: `{ products: [{ id, name, description, price, image_path, stock, created_at, updated_at }] }`
- `GET /api/products/:id` — response: `{ product: { id, name, description, price, image_path, stock, created_at, updated_at } }`
- `POST /api/products` — request: `{ name, description?, price, image_path?, stock }`; response: `{ product: { id, ... } }`

**Cart endpoints:**

- `GET /api/cart` — response: `{ items: [{ product_id, name, price, quantity }] }`
- `POST /api/cart` — request: `{ product_id, quantity }` (quantity min 1); response: `{ item: { product_id, quantity } }`
- `PUT /api/cart/:productId` — request: `{ quantity }` (quantity min 1); response: `{ item: { product_id, quantity } }`
- `DELETE /api/cart/:productId` — response: `{ message }`

**Error format:** All endpoints return `{ error: "message" }` on failure. The HTTP status code reflects the error type (400, 401, 404, 500).

The API base URL defaults to `http://localhost:8080` and can be overridden with the `REACT_APP_API_URL` environment variable (CRA convention).

The `src/services/api.ts` module provides a thin `fetch` wrapper that:
- Prepends the base URL
- Attaches `Authorization: Bearer <token>` from localStorage when available
- Parses JSON responses
- On 401, clears the token and redirects to `/login`

### Auth State Management

A simple React context (`useAuth`) holds:
- `token: string | null` (from localStorage)
- `user: { id } | null` (populated from session.user_id after login)
- `login(email, password)`, `register(email, password)`, `logout()` actions

On app mount, if a token exists in localStorage, the app does not re-validate it proactively (the API returns 401 on expired/invalid tokens, which triggers the redirect). This avoids an extra request on every page load.

### Admin vs Client Side

The `/admin` route prefix renders admin-specific pages. Admin status is inferred from the API response — for the initial implementation, any authenticated user can access admin routes. A future iteration can add a role field to the user model for proper admin gating. Product creation (which the server currently exposes without admin check) is placed behind the admin UI to establish the convention.

### Styling

Use a single `src/index.css` or `src/App.css` for all styles. No CSS framework dependency — keep it simple with plain CSS that mimics a basic ecommerce look. The spec for the frontend repo structure (issue #2) anticipated subdirectory splitting, but the initial implementation keeps CSS flat.

### Convention Alignment

Follow the existing project conventions:
- Components in `src/components/`, pages in `src/pages/`, hooks in `src/hooks/`, services in `src/services/`
- Domain subdirectories when a domain grows (e.g., `src/services/auth/`)
- TypeScript with strict mode enabled (already in tsconfig.json)
- Functional components with hooks, no class components

## Risks

| Risk | Mitigation |
|------|------------|
| No admin role on the backend means any authenticated user can reach admin pages | Document as future work; the admin routes exist as a UI convention ready for backend role enforcement |
| react-router-dom is not yet in package.json | Add as the only new runtime dependency; test that react-scripts builds with it |
| CORS issues when client (port 3000) calls server (port 8080) | Add CORS middleware on the Go server if needed; document the proxy workaround via `package.json` proxy field |
| Cart page state can drift from server state | Fetch cart fresh on each navigation to `/cart`; optimistic updates are out of scope |
| Large product catalogs will need pagination | The initial implementation loads all products; pagination can be added later when the `ListProducts` handler supports limit/offset |

## Validation

- [ ] `react-router-dom` is added to `package.json` and `npm install` succeeds.
- [ ] `npm start` launches the dev server without errors.
- [ ] Navigating to `/` renders a product listing grid.
- [ ] Navigating to `/login` renders a login form; submitting correct credentials stores a token and redirects to home.
- [ ] Navigating to `/register` renders a registration form; successful registration redirects to login.
- [ ] Navigating to `/cart` without auth redirects to `/login`.
- [ ] After login, `/cart` displays the user's cart with add/update/remove functionality.
- [ ] `/admin` and `/admin/products` render the admin dashboard and product management pages (behind auth guard).
- [ ] Admin product creation form submits to `POST /api/products` and the new product appears in the listing.
- [ ] Logout clears the token from localStorage and redirects to home.
- [ ] All pages handle API errors gracefully (display error messages, not blank screens).
- [ ] `npm run build` completes without TypeScript or bundling errors.
