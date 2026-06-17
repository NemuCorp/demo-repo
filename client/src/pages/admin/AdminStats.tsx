import React, { useEffect, useState } from 'react';
import { getDashboardMetrics, DashboardMetrics } from '../../services/api';

interface AdminStatsProps {
  detailed?: boolean;
}

function eventLabel(type: string): string {
  const labels: Record<string, string> = {
    page_view: 'Page view',
    product_view: 'Product viewed',
    cart_add: 'Added to cart',
    cart_remove: 'Removed from cart',
    registration: 'New registration',
    login: 'Login',
    product_created: 'Product created',
  };
  return labels[type] || type;
}

function formatTime(iso: string): string {
  const d = new Date(iso);
  return d.toLocaleString();
}

function AdminStats({ detailed }: AdminStatsProps) {
  const [metrics, setMetrics] = useState<DashboardMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    getDashboardMetrics()
      .then((data) => {
        if (!cancelled) {
          setMetrics(data);
          setLoading(false);
        }
      })
      .catch((err) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : 'Failed to load metrics');
          setLoading(false);
        }
      });
    return () => { cancelled = true; };
  }, []);

  if (loading) {
    return <div className="loading">Loading dashboard...</div>;
  }

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (!metrics) {
    return null;
  }

  return (
    <div className="analytics">
      <div className="analytics-grid">
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.total_users}</div>
          <div className="analytics-card-label">Total Users</div>
        </div>
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.total_products}</div>
          <div className="analytics-card-label">Total Products</div>
        </div>
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.page_views}</div>
          <div className="analytics-card-label">Page Views</div>
        </div>
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.product_views}</div>
          <div className="analytics-card-label">Product Views</div>
        </div>
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.cart_adds}</div>
          <div className="analytics-card-label">Cart Adds</div>
        </div>
        <div className="analytics-card">
          <div className="analytics-card-value">{metrics.registrations}</div>
          <div className="analytics-card-label">Registrations</div>
        </div>
      </div>

      {metrics.today && (
        <div className="analytics-section">
          <h3>Today</h3>
          <div className="analytics-grid analytics-grid-sm">
            <div className="analytics-card analytics-card-sm">
              <div className="analytics-card-value">{metrics.today.active_users}</div>
              <div className="analytics-card-label">Active Users</div>
            </div>
            <div className="analytics-card analytics-card-sm">
              <div className="analytics-card-value">{metrics.today.total_events}</div>
              <div className="analytics-card-label">Events</div>
            </div>
            <div className="analytics-card analytics-card-sm">
              <div className="analytics-card-value">{metrics.today.products_viewed}</div>
              <div className="analytics-card-label">Products Viewed</div>
            </div>
          </div>
        </div>
      )}

      {metrics.top_products && metrics.top_products.length > 0 && (
        <div className="analytics-section">
          <h3>Top Products</h3>
          <table className="admin-table analytics-table">
            <thead>
              <tr>
                <th>Rank</th>
                <th>Product</th>
                <th>Views</th>
              </tr>
            </thead>
            <tbody>
              {metrics.top_products.map((tp, i) => (
                <tr key={tp.product_id}>
                  <td>{i + 1}</td>
                  <td>{tp.product_name}</td>
                  <td>{tp.views}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {detailed && metrics.recent_activity && metrics.recent_activity.length > 0 && (
        <div className="analytics-section">
          <h3>Recent Activity</h3>
          <table className="admin-table analytics-table">
            <thead>
              <tr>
                <th>Time</th>
                <th>User</th>
                <th>Event</th>
                <th>Details</th>
              </tr>
            </thead>
            <tbody>
              {metrics.recent_activity.map((ev) => (
                <tr key={ev.id}>
                  <td className="activity-time">{formatTime(ev.created_at)}</td>
                  <td>{ev.user_email}</td>
                  <td>
                    <span className={`event-badge event-${ev.event_type}`}>
                      {eventLabel(ev.event_type)}
                    </span>
                  </td>
                  <td className="activity-data">
                    {ev.event_data && typeof ev.event_data === 'object'
                      ? (ev.event_data as Record<string, string>).product_name || ''
                      : ''}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {detailed && metrics.daily_stats && metrics.daily_stats.length > 0 && (
        <div className="analytics-section">
          <h3>Daily Stats (Last 7 Days)</h3>
          <table className="admin-table analytics-table">
            <thead>
              <tr>
                <th>Day</th>
                <th>Event</th>
                <th>Count</th>
              </tr>
            </thead>
            <tbody>
              {metrics.daily_stats.map((ds, i) => (
                <tr key={`${ds.day}-${ds.event_type}-${i}`}>
                  <td>{ds.day}</td>
                  <td>
                    <span className={`event-badge event-${ds.event_type}`}>
                      {eventLabel(ds.event_type)}
                    </span>
                  </td>
                  <td>{ds.count}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

export default AdminStats;
