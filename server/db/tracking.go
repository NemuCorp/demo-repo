package db

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/NemuCorp/demo-repo/server/logger"
)

type TrackingDB struct {
	recordEvent       *sql.Stmt
	getEventsByUser   *sql.Stmt
	getRecentEvents   *sql.Stmt
	getTotalUsers     *sql.Stmt
	getTotalProducts  *sql.Stmt
	getEventCount     *sql.Stmt
	getTopProducts    *sql.Stmt
	getDailyStats     *sql.Stmt
	getTodayStats     *sql.Stmt
}

func NewTrackingDB(conn *sql.DB) (*TrackingDB, error) {
	var t TrackingDB

	stmt, err := conn.Prepare(`
		INSERT INTO analytics_events (user_id, event_type, event_data, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, event_type, event_data, created_at
	`)
	if err != nil {
		return nil, err
	}
	t.recordEvent = stmt

	stmt, err = conn.Prepare(`
		SELECT id, user_id, event_type, event_data, created_at
		FROM analytics_events
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`)
	if err != nil {
		return nil, err
	}
	t.getEventsByUser = stmt

	stmt, err = conn.Prepare(`
		SELECT ae.id, ae.user_id, ae.event_type, ae.event_data, ae.created_at,
			   COALESCE(u.email, 'anonymous') as user_email
		FROM analytics_events ae
		LEFT JOIN users u ON u.id = ae.user_id
		ORDER BY ae.created_at DESC
		LIMIT $1
	`)
	if err != nil {
		return nil, err
	}
	t.getRecentEvents = stmt

	stmt, err = conn.Prepare(`SELECT COUNT(*) FROM users`)
	if err != nil {
		return nil, err
	}
	t.getTotalUsers = stmt

	stmt, err = conn.Prepare(`SELECT COUNT(*) FROM products`)
	if err != nil {
		return nil, err
	}
	t.getTotalProducts = stmt

	stmt, err = conn.Prepare(`
		SELECT COUNT(*) FROM analytics_events WHERE event_type = $1
	`)
	if err != nil {
		return nil, err
	}
	t.getEventCount = stmt

	stmt, err = conn.Prepare(`
		SELECT event_data->>'product_id' as product_id,
			   event_data->>'product_name' as product_name,
			   COUNT(*) as views
		FROM analytics_events
		WHERE event_type = 'product_view'
		GROUP BY event_data->>'product_id', event_data->>'product_name'
		ORDER BY views DESC
		LIMIT $1
	`)
	if err != nil {
		return nil, err
	}
	t.getTopProducts = stmt

	stmt, err = conn.Prepare(`
		SELECT DATE(created_at) as day, event_type, COUNT(*) as count
		FROM analytics_events
		WHERE created_at >= $1
		GROUP BY DATE(created_at), event_type
		ORDER BY day DESC
		LIMIT $2
	`)
	if err != nil {
		return nil, err
	}
	t.getDailyStats = stmt

	stmt, err = conn.Prepare(`
		SELECT COUNT(DISTINCT user_id) as active_users,
			   COUNT(*) as total_events,
			   COUNT(DISTINCT CASE WHEN event_type = 'product_view' THEN event_data->>'product_id' END) as products_viewed
		FROM analytics_events
		WHERE created_at >= $1
	`)
	if err != nil {
		return nil, err
	}
	t.getTodayStats = stmt

	logger.Info.Println("TrackingDB prepared statements initialized")
	return &t, nil
}

type AnalyticsEvent struct {
	ID        int              `json:"id"`
	UserID    *int             `json:"user_id"`
	EventType string           `json:"event_type"`
	EventData *json.RawMessage `json:"event_data"`
	CreatedAt time.Time        `json:"created_at"`
}

type RecentEvent struct {
	ID        int              `json:"id"`
	UserID    *int             `json:"user_id"`
	UserEmail string           `json:"user_email"`
	EventType string           `json:"event_type"`
	EventData *json.RawMessage `json:"event_data"`
	CreatedAt time.Time        `json:"created_at"`
}

type TopProduct struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Views       int    `json:"views"`
}

type DailyStat struct {
	Day       string `json:"day"`
	EventType string `json:"event_type"`
	Count     int    `json:"count"`
}

type TodayStats struct {
	ActiveUsers    int `json:"active_users"`
	TotalEvents    int `json:"total_events"`
	ProductsViewed int `json:"products_viewed"`
}

type DashboardMetrics struct {
	TotalUsers      int          `json:"total_users"`
	TotalProducts   int          `json:"total_products"`
	PageViews       int          `json:"page_views"`
	ProductViews    int          `json:"product_views"`
	CartAdds        int          `json:"cart_adds"`
	Registrations   int          `json:"registrations"`
	Today           *TodayStats  `json:"today"`
	TopProducts     []TopProduct `json:"top_products"`
	RecentActivity  []RecentEvent `json:"recent_activity"`
	DailyStats      []DailyStat  `json:"daily_stats"`
}

func (t *TrackingDB) RecordEvent(userID *int, eventType string, eventData map[string]interface{}, ipAddress, userAgent string) (*AnalyticsEvent, error) {
	dataJSON, err := json.Marshal(eventData)
	if err != nil {
		return nil, err
	}

	ae := &AnalyticsEvent{}
	var raw json.RawMessage
	err = t.recordEvent.QueryRow(userID, eventType, dataJSON, ipAddress, userAgent).Scan(
		&ae.ID, &ae.UserID, &ae.EventType, &raw, &ae.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	ae.EventData = &raw
	return ae, nil
}

func (t *TrackingDB) GetRecentEvents(limit int) ([]RecentEvent, error) {
	rows, err := t.getRecentEvents.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []RecentEvent
	for rows.Next() {
		var e RecentEvent
		var raw json.RawMessage
		if err := rows.Scan(&e.ID, &e.UserID, &e.EventType, &raw, &e.CreatedAt, &e.UserEmail); err != nil {
			return nil, err
		}
		e.EventData = &raw
		events = append(events, e)
	}
	if events == nil {
		events = []RecentEvent{}
	}
	return events, rows.Err()
}

func (t *TrackingDB) GetDashboardMetrics() (*DashboardMetrics, error) {
	m := &DashboardMetrics{}

	t.getTotalUsers.QueryRow().Scan(&m.TotalUsers)
	t.getTotalProducts.QueryRow().Scan(&m.TotalProducts)
	t.getEventCount.QueryRow("page_view").Scan(&m.PageViews)
	t.getEventCount.QueryRow("product_view").Scan(&m.ProductViews)
	t.getEventCount.QueryRow("cart_add").Scan(&m.CartAdds)
	t.getEventCount.QueryRow("registration").Scan(&m.Registrations)

	m.Today = &TodayStats{}
	todayStart := time.Now().Truncate(24 * time.Hour)
	t.getTodayStats.QueryRow(todayStart).Scan(&m.Today.ActiveUsers, &m.Today.TotalEvents, &m.Today.ProductsViewed)

	m.TopProducts = []TopProduct{}
	rows, err := t.getTopProducts.Query(5)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var tp TopProduct
			if err := rows.Scan(&tp.ProductID, &tp.ProductName, &tp.Views); err != nil {
				continue
			}
			m.TopProducts = append(m.TopProducts, tp)
		}
	}

	m.RecentActivity, _ = t.GetRecentEvents(20)

	m.DailyStats = []DailyStat{}
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	dRows, err := t.getDailyStats.Query(sevenDaysAgo, 100)
	if err == nil {
		defer dRows.Close()
		for dRows.Next() {
			var ds DailyStat
			if err := dRows.Scan(&ds.Day, &ds.EventType, &ds.Count); err != nil {
				continue
			}
			m.DailyStats = append(m.DailyStats, ds)
		}
	}

	return m, nil
}
