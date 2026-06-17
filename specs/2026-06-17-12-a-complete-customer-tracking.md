# A Complete Customer Tracking Solution

**Issue:** [NemuCorp/demo-repo#12](https://github.com/NemuCorp/demo-repo/issues/12)

## Problem

The codebase currently has zero customer tracking. The only user-related data stored are account credentials (`users` table) and auth sessions (`sessions` table). There is no visibility into how customers use the ecommerce site — what pages they visit, which products they view, how they navigate the purchase funnel, or where they abandon the cart. The admin dashboard is a placeholder that displays no real metrics.

Without tracking, the business cannot answer basic questions: Who are our most active users? Which products drive the most interest? Where are users dropping off? What is the conversion rate from product view to cart-add to purchase?

## Goals

- **Design a minimal, privacy-respecting tracking system** that captures customer behavior within the ecommerce site — not third-party tracking, no cross-site cookies, no fingerprinting.
- **Track the right metrics:** page views, product impressions, cart additions, cart removals, and session duration. These are actionable business metrics (not raw clickstreams).
- **Do NOT track:** IP addresses, PII beyond the existing user email, keystrokes, cursor movements, or off-site behavior.
- **Store tracking data in PostgreSQL** alongside existing domain tables, keeping the architecture simple (no new data stores).
- **Expose analytics endpoints** in the Go backend that aggregate tracking data for the admin dashboard.
- **Build a real admin analytics dashboard** that replaces the placeholder with live charts and metrics, providing genuine feedback to the operator.
- **Auto-instrument the frontend** with a lightweight tracking hook/page-view logger that requires no per-page wiring beyond a single provider wrapper.
- **Keep dependencies minimal:** add a charting library (recharts) to the client for admin visualizations; no new backend dependencies beyond the standard library.

## Approach

### Scenarios and Use Cases

#### Scenario 1 — Anonymous Browsing

A visitor lands on the home page without logging in. They browse the product listing, click into a product detail page, and view several products. The system records each page view as an anonymous event tied to a browser-side anonymous session ID (a UUID generated on first visit and stored in `localStorage`). When the visitor later registers or logs in, the anonymous events are stitched to the authenticated user by replacing the anonymous session ID with the user ID.

**Metrics captured:** page views (home, product listing, product detail), product impressions (per-product view count), time-on-page (via heartbeat pings), session start/end.

#### Scenario 2 — Authenticated Shopping

A logged-in user browses products, adds items to their cart, updates quantities, and removes items. Each cart action is recorded as an event. The user's full journey — from landing through cart to (future) checkout — is traceable in the admin dashboard.

**Metrics captured:** all page views (now attributed to a known user), cart-add events (product_id, quantity), cart-remove events, cart-update events (old quantity → new quantity), session duration.

#### Scenario 3 — Admin Analytics Review

An admin opens the analytics dashboard. They see:
- **Overview cards:** total users, active users today, new sign-ups this week, total page views today.
- **Time-series chart:** page views over the last 30 days (line chart).
- **Top products:** a ranked table of products by impression count and cart-add count.
- **User activity table:** a searchable, sortable list of recent user sessions with page view counts and last activity time.
- **Conversion funnel:** product views → cart adds → (placeholder for purchases).

All data is near-real-time — tracking events are written synchronously (non-blocking to the main response) and aggregated on read.

#### Scenario 4 — Privacy and Retention

Tracking data older than 90 days is automatically pruned (matching the 90-day session expiration pattern). Users can request deletion of their tracking data through a future account page. The tracking system does not record IP addresses, User-Agent strings, or referrer headers — only the in-app page URL, product ID where applicable, timestamp, and the user/session identifier.

### What Is Tracked vs. What Is Not

| Data Point | Tracked? | Reason |
|---|---|---|
| Page URL visited | Yes | Core analytics: which pages are popular |
| Product ID viewed | Yes | Product interest/impression metrics |
| Cart actions (add/update/remove) | Yes | Purchase funnel analysis |
| Timestamp of event | Yes | Time-series analysis |
| User ID (if authenticated) | Yes | Per-user analytics |
| Anonymous session ID | Yes | Pre-login funnel attribution |
| Session duration (start/end) | Yes | Engagement metrics |
| IP address | No | Not needed; privacy risk |
| User-Agent / browser info | No | Not actionable for this use case |
| Referrer header | No | No external traffic sources to analyze |
| Keystrokes / cursor movements | No | Excessive; privacy-invasive |
| Password or session token | No | Security risk |
| Cart contents snapshot | No | Cart state already stored in `cart_items` |
| Off-site behavior | No | Out of scope |

### Database Schema

Three new tables in a new migration file (`server/db/migrations/002_tracking.sql`):

```sql
-- Anonymous sessions for pre-auth tracking
CREATE TABLE IF NOT EXISTS anonymous_sessions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     INTEGER REFERENCES users(id) ON DELETE SET NULL,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_seen   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Individual tracking events
CREATE TABLE IF NOT EXISTS tracking_events (
    id              BIGSERIAL PRIMARY KEY,
    anonymous_id    UUID REFERENCES anonymous_sessions(id) ON DELETE SET NULL,
    user_id         INTEGER REFERENCES users(id) ON DELETE SET NULL,
    event_type      VARCHAR(50) NOT NULL,  -- page_view, product_view, cart_add, cart_remove, cart_update, session_start, session_end
    page_url        VARCHAR(500),
    product_id      INTEGER REFERENCES products(id) ON DELETE SET NULL,
    metadata        JSONB DEFAULT '{}',
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Indexes for common analytics queries
CREATE INDEX idx_tracking_events_user_id ON tracking_events(user_id);
CREATE INDEX idx_tracking_events_anonymous_id ON tracking_events(anonymous_id);
CREATE INDEX idx_tracking_events_event_type ON tracking_events(event_type);
CREATE INDEX idx_tracking_events_created_at ON tracking_events(created_at);
CREATE INDEX idx_tracking_events_product_id ON tracking_events(product_id);
CREATE INDEX idx_tracking_events_user_time ON tracking_events(user_id, created_at DESC);
```

The `metadata` JSONB column stores event-specific details (e.g., `{"quantity": 3}` for cart_add, `{"prev_quantity": 2, "new_quantity": 5}` for cart_update) without requiring schema changes per event type.

### Backend Changes

#### New DB Layer (`server/db/tracking.go`)

A `TrackingDB` struct with prepared statements:
- `insertEvent(anonymousID, userID, eventType, pageURL, productID, metadata)` — insert a single event.
- `getEventsByUser(userID, limit, offset)` — paginated event list for a user.
- `getAggregateStats()` — returns: total users, active users today, new users this week, total events today, total events last 30 days.
- `getDailyEventCounts(days)` — array of `{date, count}` for time-series charts.
- `getTopProductsByImpressions(limit)` — top products ranked by product_view events.
- `getTopProductsByCartAdds(limit)` — top products ranked by cart_add events.
- `getRecentUserSessions(limit, offset)` — recent anonymous_sessions rows with page-view counts.
- `pruneOldEvents(retentionDays)` — deletes events older than retention period.
- `createAnonymousSession()` — creates a new anonymous session, returns the UUID.
- `linkAnonymousSession(anonymousID, userID)` — associates an anonymous session with a user after login/register.

#### New Handler Layer (`server/handler/analytics.go`)

| Method | Path | Auth | Description |
|--------|------|:---:|-------------|
| POST | `/api/analytics/event` | No | Record a tracking event (from frontend auto-tracker). Accepts body `{ anonymous_id?, event_type, page_url?, product_id?, metadata? }`. If a valid auth token is present, user_id is extracted from the session and attached. |
| GET | `/api/analytics/stats` | Yes | Aggregate stats for the admin dashboard: total users, active today, new this week, events today, events last 30 days. |
| GET | `/api/analytics/events/daily?days=30` | Yes | Time-series page-view counts per day for charting. |
| GET | `/api/analytics/products/top?limit=10&by=impressions` | Yes | Top products ranked by impressions or cart_adds. |
| GET | `/api/analytics/users/recent?limit=20&offset=0` | Yes | Recent user sessions with activity summaries. |

#### Tracking Middleware (`server/handler/middleware.go`)

Extend the existing middleware or add a lightweight tracking middleware that:
- Generates or extracts an anonymous session ID from a `X-Anonymous-ID` header (sent by the frontend from localStorage).
- Optionally records page_view events on request (can be deferred to the client-side tracker for cleaner separation; the middleware approach is documented as an alternative for server-side accuracy).
- Does NOT block or delay the response on tracking failures (event insert failures are logged and swallowed).

#### Anonymous Session ID Flow

1. On first page load, the frontend generates a UUID and stores it in `localStorage` under key `anonymous_session_id`.
2. The frontend sends this UUID in the `X-Anonymous-ID` header on API calls or as part of the event payload.
3. When the user registers or logs in, the frontend receives the user ID and calls `POST /api/analytics/link` (or the backend links automatically during login/register) to associate the anonymous session with the user.
4. Future events carry both the (optional) anonymous ID and the authenticated user ID.

### Frontend Changes

#### Tracking Infrastructure (`client/src/tracking/`)

```
src/tracking/
  AnalyticsProvider.tsx    # React context wrapping the app; initializes anonymous session ID, provides `track()` function
  tracker.ts              # Core tracking logic: generate/manage anonymous session ID, batch-send events to POST /api/analytics/event
  usePageView.ts           # Hook that auto-fires a page_view event on route change
  useProductView.ts        # Hook that fires a product_view event when a product detail page mounts
```

**`AnalyticsProvider`** wraps the entire `<App />` (or `<BrowserRouter>`) and:
- On mount, checks `localStorage` for `anonymous_session_id`; generates one if absent.
- Provides a `track(eventType, payload)` function to the React tree via context.
- Batching: accumulates events in a queue and flushes every 5 seconds or when the queue reaches 10 events (whichever comes first) to reduce HTTP overhead.
- On auth state change (user logs in or out), updates the user ID on future events.
- On unmount (page close/navigation away), flushes any remaining events via `navigator.sendBeacon` for reliable delivery.

**`usePageView`** uses `useLocation()` from react-router-dom to fire a `page_view` event on every route change. It also starts/ends a session heartbeat (a periodic ping every 30 seconds while the tab is active, tracked via `document.visibilitychange`).

**`useProductView`** fires a `product_view` event when a product detail page mounts, including the `product_id` in the metadata.

#### Cart Event Instrumentation

The existing cart service functions (`addToCart`, `updateCartItem`, `removeFromCart`) are extended to also call `track()` with `cart_add`, `cart_update`, or `cart_remove` event types and relevant metadata (product_id, quantity). This ties cart actions into the tracking pipeline without changing the cart API contract.

#### Admin Analytics Dashboard (`client/src/pages/admin/Analytics.tsx`)

Replaces the placeholder admin dashboard with a rich analytics page. Depends on `recharts` (added to `package.json`) for charts.

Layout (top to bottom):

1. **Stats Overview Row** — Four metric cards (total users, active today, new sign-ups this week, total page views today), each with a label, large number, and a subtle up/down trend indicator.

2. **Page Views Over Time** — A line chart (`<LineChart>` from recharts) showing daily page-view counts for the last 30 days. The admin can toggle between 7, 30, and 90-day windows.

3. **Top Products** — Two-column layout: left column shows "Top by Impressions" (ranked table of products by product_view count), right column shows "Top by Cart Adds" (ranked table by cart_add count). Each row shows product name, count, and a mini sparkline or bar.

4. **Recent User Activity** — A paginated table listing recent users/sessions with: email (or "Anonymous"), last seen time, page view count, and a "View Details" button that expands an inline activity log (last 20 events for that user).

Fetching uses the new analytics API endpoints. The page refreshes automatically every 60 seconds (with a manual refresh button).

#### New Admin Route

```
/admin/analytics     → Analytics Dashboard (replaces the placeholder /admin landing)
/admin/products      → Product Management (existing, unchanged)
```

The `/admin` route is updated to default-redirect to `/admin/analytics`. The admin nav (`AdminDashboard.tsx`) is updated with an "Analytics" link alongside the existing "Manage Products" link.

### Dependency Changes

**Client (`client/package.json`):**
- Add `recharts` (v2.x) for admin charts.

**Server:** No new dependencies. Use existing `database/sql` for new tables.

## Risks

| Risk | Mitigation |
|------|------------|
| Event table grows unboundedly large | Implement automatic pruning of events older than 90 days via a periodic cleanup goroutine or on-read filtering. Add a `clean` CLI subcommand entry for manual truncation. |
| Frontend tracking adds latency to page loads | Events are sent asynchronously; the `track()` function uses `navigator.sendBeacon` or fires-and-forgets. No page render is blocked on tracking I/O. |
| Anonymous-to-user stitching is complex and error-prone | Keep it simple: on login/register, update the `anonymous_sessions` row's `user_id` via `linkAnonymousSession()`. Past events retain their anonymous ID; queries combine events by user_id OR anonymous_id linked to that user. |
| recharts adds a non-trivial dependency to the bundle | recharts is tree-shakeable; only the used chart components (`LineChart`, `BarChart`, `PieChart`) are imported. The analytics page is code-split via `React.lazy()` so it does not affect the public-facing pages' bundle size. |
| Database migration 002 may conflict with future changes | Number migrations sequentially. The current migration is `001_initial.sql`; `002_tracking.sql` follows naturally. |
| High event volume under load | Events are insert-only with minimal indexed columns. PostgreSQL handles this well for the expected demo scale. For production scale, a future iteration could introduce an async queue (e.g., Redis + worker) for event ingestion. |
| No admin role enforcement (existing gap) | The analytics endpoints use auth middleware (same as cart). Any authenticated user can see analytics. This matches the existing pattern; proper admin role gating is tracked separately. |

## Validation

- [ ] Migration `002_tracking.sql` creates all three tables and indexes without error.
- [ ] `go build ./...` in `server/` succeeds with new `tracking.go` DB layer and `analytics.go` handler.
- [ ] `POST /api/analytics/event` accepts and stores a tracking event; returns 200 with `{}` on success, 400 on invalid payload.
- [ ] `GET /api/analytics/stats` returns correct aggregate counts (users, active today, events today, etc.).
- [ ] `GET /api/analytics/events/daily?days=30` returns an array of `{date, count}` objects for charting.
- [ ] `GET /api/analytics/products/top?by=impressions` returns top products ranked by view count.
- [ ] `GET /api/analytics/users/recent` returns paginated recent sessions with activity summaries.
- [ ] Frontend generates and persists an `anonymous_session_id` in `localStorage` on first visit.
- [ ] Navigating between pages in the client fires `page_view` events sent to `POST /api/analytics/event`.
- [ ] Viewing a product detail page fires a `product_view` event with the correct product_id.
- [ ] Adding/updating/removing cart items fires `cart_add`/`cart_update`/`cart_remove` events.
- [ ] After login, the anonymous session is linked to the user and subsequent events carry the user_id.
- [ ] The admin analytics page (`/admin/analytics`) renders without errors and displays live data.
- [ ] Stats cards show correct total users, active today, new sign-ups, and page views.
- [ ] The line chart renders daily page-view data and responds to time-window toggles (7/30/90 days).
- [ ] Top products tables show correct rankings for both impressions and cart adds.
- [ ] Recent user activity table is populated and sortable.
- [ ] `npm run build` in `client/` completes without TypeScript or bundling errors.
- [ ] Existing tests (if any) continue to pass.
