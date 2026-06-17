# Docker Compose Up Resulting in Error

**Issue:** [NemuCorp/demo-repo#27](https://github.com/NemuCorp/demo-repo/issues/27)

## Problem

The `docker compose up --build` command fails during the client build stage. The `npm run build` step (`react-scripts build`) emits a TypeScript compilation error:

```
TS2459: Module '"../services/api"' declares 'Product' locally, but it is not exported.
```

This error breaks the multi-stage Docker build for the `client` service (`client/Dockerfile` line 6: `RUN npm run build`), preventing `docker compose up` from completing.

### Affected files

| File | Import | Problem |
|------|--------|---------|
| `client/src/components/ProductCard.tsx:3` | `import { Product } from '../services/api'` | TS2459 — `Product` not exported |
| `client/src/pages/ProductPage.tsx:3` | `import { Product, getProduct, addToCart } from '../services/api'` | TS2459 — `Product` not exported |
| `client/src/pages/HomePage.tsx:2` | `import { Product, getProducts } from '../services/api'` | TS2459 — `Product` not exported |
| `client/src/pages/admin/AdminProducts.tsx` | `import { Product, getProducts, createProduct } from '../../services/api'` | TS2459 — `Product` not exported |
| `client/src/pages/CartPage.tsx` | `import { CartItem, getCart, updateCartItem, removeFromCart } from '../services/api'` | `removeFromCart` does not exist — exported name is `removeCartItem` |

### Root cause

`client/src/services/api.ts` imports `Product` and `CartItem` from `../types` and uses them internally in function return types, but never re-exports these types. Any consumer file that attempts a named `import { Product }` or `import { CartItem }` from `../services/api` hits TS2459 because the module declares them locally without exporting them.

The `CartPage.tsx` import also references `removeFromCart`, but the actual function exported from `api.ts` is named `removeCartItem`.

## Goals

- Fix all TypeScript compilation errors so `npm run build` succeeds both locally and inside the Docker client build stage.
- Restore the ability to run `docker compose up --build` without errors.
- Ensure zero TypeScript errors from `react-scripts build` across the entire client app.

## Approach

### Primary fix: Re-export `Product` and `CartItem` from `api.ts`

Add type-only re-exports in `client/src/services/api.ts`:

```typescript
import { AuthResponse, CartItem, Product, User } from '../types';

// ... existing code ...

export type { Product };
export type { CartItem };
```

This is a one-file, two-line change that fixes all four consumer files. It keeps existing import structures untouched and is consistent with the precedent that `api.ts` already exports its own `DashboardMetrics` interface inline.

### Secondary fix: Fix `removeFromCart` naming mismatch in `CartPage.tsx`

Change the import in `client/src/pages/CartPage.tsx` from `removeFromCart` to `removeCartItem` to match the actual exported function name, and update the corresponding usage site.

### Scope

- Modify `client/src/services/api.ts` — add `export type { Product }` and `export type { CartItem }`.
- Modify `client/src/pages/CartPage.tsx` — rename `removeFromCart` import and usage to `removeCartItem`.
- No changes to Dockerfiles, docker-compose.yml, or backend code.
- No new dependencies.

## Risks

| Risk | Mitigation |
|------|------------|
| Other consumer files may have similar broken type imports | Audit all imports from `../services/api` during validation |
| The `CartPage.tsx` naming mismatch may not surface until `CartItem` re-export also exposes the `removeFromCart` error | Fix both issues in one pass; validate with `npm run build` |
| `Product` is already re-exported from `api.ts` but `CartItem` consumers may still break | Include `CartItem` in the re-export list proactively |

## Validation

- [ ] `export type { Product }` and `export type { CartItem }` added to `client/src/services/api.ts`.
- [ ] `removeFromCart` renamed to `removeCartItem` in `client/src/pages/CartPage.tsx`.
- [ ] `npm run build` completes with zero errors in the `client/` directory.
- [ ] `docker compose up --build` succeeds — all three services (db, server, client) start without build errors.
- [ ] No new TypeScript or lint warnings introduced.
