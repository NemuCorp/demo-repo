# Create a Suitable Frontend

**Issue:** [NemuCorp/demo-repo#7](https://github.com/NemuCorp/demo-repo/issues/7)

## Problem

The `client/` directory contains a bare placeholder: `App.tsx` renders "Demo Repo — Client application is running" and `index.tsx` mounts it to the DOM. There is no routing, no pages, and no integration with the backend API. The Go server (`server/`) already exposes a full ecommerce API — auth (register/login/logout), products (list/get/create), and cart (view/add/update/remove) — and the React frontend must consume it. The UI must serve two audiences: public shoppers browsing products and managing their cart, and authenticated administrators managing the product catalog.

## Goals

- Install `react-router-dom` v6 for client-side routing and split the app into public (shopper) and admin routes.
- Build reusable UI components: navigation header, product cards, cart item rows, and form controls.
- Create public-facing pages: Home (product listing), Product Detail, Cart, Login, and Register.
- Create admin-facing pages: Dashboard and Product Management (list and create).
- Build a typed API service layer under `src/services/` to call the backend endpoints, with auth token management.
- Implement auth state via a React context hook that persists the Bearer token in `localStorage`.
- Guard protected routes (cart, admin) with a wrapper that redirects unauthenticated users to login.
- Keep dependencies minimal: only `react-router-dom` as a new runtime dependency; use native `fetch` and React state for all other functionality.

## Approach

### Dependencies

Add `react-router-dom` (v6) to `client/package.json`. No other new runtime dependencies. The project already uses `react-scripts` (CRA), TypeScript with strict mode, and React 18 — all tooling stays as-is.

### Route Structure

```
/                    → Home (public product listing grid)
/products/:id        → Product Detail (single product view with add-to-cart)
/cart                → Cart (auth required; view/edit/remove items)
/login               → Login form
/register            → Registration form
/admin               → Admin Dashboard (auth required; placeholder summary)
/admin/products      → Admin Product Management (list + create; auth required)
```

A `ProtectedRoute` wrapper component gates `/cart`, `/admin`, and `/admin/products`: if no auth token exists in `localStorage`, it redirects to `/login`.

### Component Tree

```
src/
  components/
    Header.tsx          # Top nav bar; links change based on auth state and admin routes
    ProductCard.tsx      # Thumbnail card used in Home grid and Product Detail
    CartItemRow.tsx      # Single cart row with quantity stepper and remove button
    ProtectedRoute.tsx   # Auth guard that redirects to /login when unauthenticated
    AdminHeader.tsx      # Admin-specific navigation (Dashboard, Products)
  pages/
    Home.tsx             # Public product grid fetched from GET /api/products
    ProductDetail.tsx    # Single product view fetched from GET /api/products/:id
    Cart.tsx             # Cart view; fetches GET /api/cart, supports update/remove
    Login.tsx            # Login form; calls POST /api/auth/login
    Register.tsx         # Registration form; calls POST /api/auth/register
    admin/
      Dashboard.tsx      # Admin landing with placeholder summary stats
      ProductList.tsx    # Admin product table (list + create via link/form)
      ProductForm.tsx    # Create-product form; calls POST /api/products
  hooks/
    useAuth.tsx          # Auth context + provider; holds token and user id
    useApi.tsx           # Thin fetch wrapper that injects the Bearer token
  services/
    api.ts               # Base fetch helper (prepend API URL, attach Authorization header, handle 401 → redirect)
    auth.ts              # login(), register(), logout()
    products.ts          # listProducts(), getProduct(), createProduct()
    cart.ts              # getCart(), addToCart(), updateCartItem(), removeCartItem()
  App.tsx                # BrowserRouter + Routes configuration
  index.tsx              # (unchanged; already renders <App />)
```

### API Integration

The Go server runs on port 8080 and exposes these endpoints:

| Method | Path | Auth Required | Description |
|--------|------|:---:|-------------|
| POST | `/api/auth/register` | No | Create account; body `{ email, password }`, password min 6 chars; returns `{ user: { id, email } }` |
| POST | `/api/auth/login` | No | Authenticate; body `{ email, password }`; returns `{ token, session: { id, user_id, created_at, expires_at } }` |
| POST | `/api/auth/logout` | Yes | Invalidate all sessions for user; returns `{ message: "logged out" }` |
| GET | `/api/products` | No | List all products; returns `{ products: [{ id, name, description?, price, image_path?, stock, created_at, updated_at }] }` |
| GET | `/api/products/:id` | No | Get single product; returns `{ product: { id, name, description?, price, image_path?, stock, created_at, updated_at } }` |
| POST | `/api/products` | No | Create product; body `{ name, description?, price, image_path?, stock }`; returns `{ product: { id, ... } }` |
| GET | `/api/cart` | Yes | View user's cart; returns `{ cart: [{ id, user_id, product_id, quantity, product_name, price, image_path?, created_at }] }` |
| POST | `/api/cart` | Yes | Add item; body `{ product_id, quantity }` (quantity min 1); returns `{ item: { id, user_id, product_id, quantity, product_name, price, image_path?, created_at } }` |
| PUT | `/api/cart/:productId` | Yes | Update quantity; body `{ quantity }` (min 0, setting to 0 removes the item); returns `{ item: { ... } }` or `{ message: "item removed" }` |
| DELETE | `/api/cart/:productId` | Yes | Remove item; returns `{ message: "item removed" }` |

All error responses use shape `{ error: "message" }`. HTTP status reflects the error type: 400 (bad request), 401 (unauthorized), 404 (not found), 409 (conflict, e.g. duplicate email), 500 (internal).

The API base URL defaults to `http://localhost:8080` and can be overridden via the `REACT_APP_API_URL` environment variable (CRA convention). Set `"proxy": "http://localhost:8080"` in `client/package.json` to avoid CORS issues during development.

The `src/services/api.ts` module provides a thin `fetch` wrapper that:
- Prepends the API base URL.
- Attaches `Authorization: Bearer <token>` from `localStorage` when present.
- Parses JSON responses.
- On HTTP 401, clears the stored token and redirects to `/login`.

### Auth State Management

A React context (`useAuth`) provides:
- `token: string | null` — stored in and loaded from `localStorage`.
- `user: { id: number } | null` — populated from the login response's `session.user_id`.
- `login(email, password)`, `register(email, password)`, `logout()` — action methods.

On app mount, if a token exists in `localStorage`, the app does not proactively validate it. The API returns 401 on expired/invalid tokens, which triggers the redirect via the `api.ts` fetch wrapper. This avoids an extra round-trip on every page load.

### Admin vs Client

The `/admin` route prefix renders admin-specific pages behind the auth guard. The current backend does not enforce an admin role (any authenticated user can call `POST /api/products`), so admin access is a UI convention for now. A future iteration can add a `role` column to the `users` table and enforce it in middleware. Product creation is placed under the admin UI to establish the expected structure.

### Styling

Use a single CSS file (e.g., `src/App.css` or `src/index.css`) with plain CSS — no framework dependency. The spec for issue #2 anticipated subdirectory splitting, but the initial implementation keeps styles flat and centralized.

### Convention Alignment

Follow the project conventions established in the initial repo structure (issue #2):
- Components in `src/components/`, pages in `src/pages/`, hooks in `src/hooks/`, services in `src/services/`.
- Domain subdirectories when a domain grows (e.g., `src/pages/admin/`).
- TypeScript with strict mode enabled (already in `tsconfig.json`).
- Functional components with hooks; no class components.
- `react-scripts` for build and dev server (already in `package.json`).

## Risks

| Risk | Mitigation |
|------|------------|
| `react-router-dom` not yet in `package.json` | Add as the sole new runtime dependency; verify `npm install` and `npm start` succeed |
| CORS when client (port 3000) calls server (port 8080) | Add `"proxy": "http://localhost:8080"` to `package.json` for development; document the Go side CORS middleware alternative |
| Cart state may drift from server state | Always fetch cart fresh on `/cart` navigation; optimistic updates are out of scope for this iteration |
| Large product catalogs will need pagination | Initial implementation loads all products; pagination can be added when the `ListProducts` handler gains limit/offset support |
| No admin role enforcement on backend | Document as future work; the admin UI routes exist as a convention ready for backend role gating |

## Validation

- [ ] `react-router-dom` is added to `package.json` and `npm install` succeeds.
- [ ] `npm start` launches the dev server without errors.
- [ ] Navigating to `/` renders a product listing grid fetched from the API.
- [ ] Navigating to `/login` renders a login form; submitting valid credentials stores a token and redirects to home.
- [ ] Navigating to `/register` renders a registration form; successful registration redirects to login.
- [ ] Navigating to `/cart` without auth redirects to `/login`.
- [ ] After login, `/cart` displays the user's cart items with add/update/remove functionality.
- [ ] Navigating to `/admin` and `/admin/products` renders the admin dashboard and product management pages (behind auth guard).
- [ ] Admin product creation form submits to `POST /api/products` and the new product appears in the listing.
- [ ] Logout clears the token from `localStorage` and redirects to home.
- [ ] All pages display inline error messages when API calls fail (no blank screens or silent failures).
- [ ] `npm run build` completes with zero TypeScript or bundling errors.
