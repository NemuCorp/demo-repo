# Create Frontend

**Issue:** [NemuCorp/demo-repo#6](https://github.com/NemuCorp/demo-repo/issues/6)

## Problem

The `server/` exposes a working REST API for auth, products, and cart operations, but the `client/` is a bare React scaffold. The application has no UI to browse products, manage a shopping cart, authenticate users, or administer the store. An ecommerce site needs a fully navigable frontend with distinct client-facing (customer) and admin-facing (store manager) experiences.

## Goals

- Build a React + TypeScript SPA that consumes the existing backend API at `/api`.
- Deliver client-side pages: Home, Product Listing, Product Detail, Cart, Login, Register.
- Deliver admin-side pages: Dashboard, Product Management.
- Implement client-side routing with protected routes for authenticated-only sections (cart, admin).
- Implement full auth flow (register, login, logout) with session token persistence.
- Provide a clean, usable ecommerce-appropriate visual design.

## Approach

### Dependencies

Add to `client/`:
- **react-router-dom** — client-side routing with `BrowserRouter`, `Routes`, `Route`.
- **axios** — HTTP client for API calls, with a configured base instance targeting the backend.
- **CSS approach** — use plain CSS modules or a lightweight utility-first approach co-located with components. Avoid heavyweight CSS frameworks unless the scope demands it.

### Source Layout

```
client/src/
  components/         # Reusable UI building blocks
    Layout.tsx        # Navbar + content area + footer shell
    Navbar.tsx        # Top navigation with auth-aware links
    ProductCard.tsx   # Product thumbnail card for listings
    CartItem.tsx      # Single row in the cart view
    ProtectedRoute.tsx# Redirects unauthenticated users to /login
  pages/
    client/
      Home.tsx        # Landing / featured products
      Products.tsx    # Full product listing / browse
      ProductDetail.tsx # Single product view with add-to-cart
      Cart.tsx        # Cart management (view, update qty, remove)
      Login.tsx       # Login form
      Register.tsx    # Registration form
    admin/
      Dashboard.tsx   # Admin overview / quick stats
      Products.tsx    # Admin product CRUD (list, create, edit, delete)
  hooks/
    useAuth.tsx       # AuthContext + provider; exposes user, login, logout, register
  services/
    api.ts            # Axios instance with base URL and auth header injection
    auth.ts           # login, register, logout API calls
    products.ts       # Product CRUD API calls
    cart.ts           # Cart API calls
  App.tsx             # Router definition: public routes, protected routes
  index.tsx           # Entry point (unchanged)
```

### Routing Plan

| Path | Page | Auth Required | Admin Required |
|------|------|:---:|:---:|
| `/` | Home | No | No |
| `/products` | Products (listing) | No | No |
| `/products/:id` | ProductDetail | No | No |
| `/login` | Login | No | No |
| `/register` | Register | No | No |
| `/cart` | Cart | Yes | No |
| `/admin` | Dashboard | Yes | Yes |
| `/admin/products` | Admin Product Management | Yes | Yes |

### Auth Flow

1. On app mount, `useAuth` checks `localStorage` for a saved session token.
2. If a token exists, the provider sets the authenticated user state.
3. Axios interceptor attaches `Authorization: Bearer <token>` to every request.
4. `ProtectedRoute` component checks auth state; redirects to `/login` if unauthenticated.
5. Admin routes additionally check for an admin indicator (e.g., user role or admin flag). Since the backend currently has no role system, admin access can be gated by a known admin email list or a localStorage flag as a temporary measure. Document this as a known gap.
6. Logout clears the token from `localStorage` and resets auth state.

### API Integration

- An axios instance in `services/api.ts` is configured with `baseURL` pointing to the backend (default `http://localhost:8080/api`).
- A request interceptor reads the session token from `localStorage` and attaches it as a Bearer header.
- A response interceptor handles 401 responses by clearing auth state and redirecting to `/login`.
- Service modules (`auth.ts`, `products.ts`, `cart.ts`) wrap API calls and return typed responses.

### Admin Side

Since the backend currently lacks admin-specific endpoints or a user-role system:
- Admin product management pages call the existing public `/api/products` endpoints (GET, POST) plus any update/delete endpoints that may be added.
- Admin dashboard shows quick stats derived from client-side state or lightweight API queries.
- Admin access is temporarily gated by a hardcoded admin email check or a manually-set localStorage flag. A follow-up issue should add proper admin roles to the backend.

## Risks

| Risk | Mitigation |
|------|------------|
| Backend has no user-role system; admin access cannot be properly enforced server-side | Gate admin routes with a frontend-only flag (admin email list or localStorage) and file a follow-up issue for backend role support |
| Product mutation endpoints (PUT, DELETE) may not exist in the current backend | Build admin product UI with create and list first; add edit/delete UI only for endpoints that exist |
| CRA (react-scripts) build tooling may feel slow or limited compared to Vite | CRA is already in place; defer migration to Vite until a future issue unless it blocks development |
| No design system or component library is present; all UI must be built from scratch | Use simple, functional CSS with CSS modules; prioritize usability over polish for the initial pass |
| Session token storage in `localStorage` is vulnerable to XSS | This is acceptable for a demo; production hardening is out of scope |

## Validation

- [ ] `npm start` in `client/` launches the React dev server without errors.
- [ ] Navigating to `/` renders the Home page with featured products loaded from the API.
- [ ] `/products` lists all products; clicking a product navigates to `/products/:id` with detail view.
- [ ] `/login` and `/register` allow user authentication; session token is persisted across page reloads.
- [ ] `/cart` is accessible only when authenticated; items can be viewed, quantity updated, and removed.
- [ ] `/admin` and `/admin/products` are gated behind authentication and an admin check.
- [ ] Admin product page allows creating new products and listing existing ones.
- [ ] Navbar reflects auth state (shows Login/Register when logged out, Cart/Logout when logged in).
- [ ] All API calls include the Bearer token when authenticated; 401 responses trigger logout.
- [ ] The app is functional against a running backend with the existing API endpoints.
