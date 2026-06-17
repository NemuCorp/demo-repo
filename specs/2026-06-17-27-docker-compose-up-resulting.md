# Docker Compose Up Resulting in Error

**Issue:** [NemuCorp/demo-repo#27](https://github.com/NemuCorp/demo-repo/issues/27)

## Problem

The `docker compose up --build` command fails during the client build stage. The `npm run build` step (`react-scripts build`) emits a TypeScript compilation error:

```
TS2459: Module '"../services/api"' declares 'Product' locally, but it is not exported.
```

This error appears in at least three files:

| File | Line | Import |
|------|------|--------|
| `client/src/components/ProductCard.tsx` | 3 | `import { Product } from '../services/api';` |
| `client/src/pages/ProductPage.tsx` | 3 | `import { Product, getProduct, addToCart } from '../services/api';` |
| `client/src/pages/HomePage.tsx` | 2 | `import { Product, getProducts } from '../services/api';` |

**Root cause:** `client/src/services/api.ts` imports `Product` from `../types` (line 1) and uses it internally in function return types (e.g., `getProducts`, `getProduct`, `createProduct`), but never re-exports the type. Consequently, any consumer file that attempts a named `import { Product }` from `../services/api` hits the TS2459 error because the module declares `Product` locally without exporting it.

The error breaks the multi-stage Docker build for the `client` service (`client/Dockerfile` line 6: `RUN npm run build`), preventing `docker compose up` from completing.

## Goals

- Fix the TypeScript compilation error so `npm run build` succeeds both locally and inside the Docker client build stage.
- Restore the ability to run `docker compose up --build` without errors.
- Ensure zero TypeScript errors from `react-scripts build` across the entire client app.

## Approach

There are two viable fixes. The chosen approach is **option A** because it is the minimal change that fixes all consumers in one place and mirrors how other shared types (`DashboardMetrics`) are already exported from the same module.

### Option A (preferred): Re-export `Product` from `api.ts`

Add a type-only re-export of `Product` in `client/src/services/api.ts`:

```typescript
import { AuthResponse, CartItem, Product, User } from '../types';

// ... existing code ...

export type { Product };  // re-export for consumers that import Product from api
```

This is a one-file, one-line change. It keeps the existing import structure in consumer files untouched and is consistent with the precedent that `api.ts` already exports its own `DashboardMetrics` interface.

### Option B (alt): Fix imports in consumer files

Change the three consumer files to import `Product` from `../types` instead of `../services/api`:

```typescript
// Before
import { Product } from '../services/api';
// After
import { Product } from '../types';
```

This approach is more correct from a dependency direction standpoint but requires changes in three files and diverges from the existing pattern where some other files also import types from `../services/api`.

### Scope

- Modify `client/src/services/api.ts` only.
- No changes to Dockerfiles, docker-compose.yml, or backend code.
- No new dependencies.

## Risks

| Risk | Mitigation |
|------|------------|
| `CartItem` may have the same TS2459 error in `CartPage.tsx` (`import { CartItem, getCart, updateCartItem, removeFromCart } from '../services/api'`) | `CartItem` should also be re-exported from `api.ts` to preempt the error; `removeFromCart` at that import site is a separate naming mismatch (`removeCartItem` vs `removeFromCart`) that may also need fixing |
| Other consumer files may have similar broken type imports | Audit all imports from `../services/api` during validation to catch any additional TS2459 errors before they surface |

## Validation

- [ ] `export type { Product };` (and `export type { CartItem };`) added to `client/src/services/api.ts`.
- [ ] `npm run build` completes with zero errors in the `client/` directory.
- [ ] `docker compose up --build` succeeds — all three services (db, server, client) start without build errors.
- [ ] No new TypeScript or lint warnings introduced.
