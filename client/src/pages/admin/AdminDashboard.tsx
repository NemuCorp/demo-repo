import React from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';

function AdminDashboard() {
  const location = useLocation();

  return (
    <div className="page">
      <h1>Admin</h1>
      <div className="admin-nav">
        <Link
          to="/admin"
          className={location.pathname === '/admin' ? 'active' : ''}
        >
          Dashboard
        </Link>
        <Link
          to="/admin/products"
          className={location.pathname === '/admin/products' ? 'active' : ''}
        >
          Manage Products
        </Link>
      </div>
      <div className="admin-content">
        {location.pathname === '/admin' ? (
          <div className="admin-welcome">
            <h2>Admin Dashboard</h2>
            <p>Welcome to the admin panel. Manage your store from here.</p>
            <div className="admin-links">
              <Link to="/admin/products" className="btn-primary">
                Manage Products
              </Link>
            </div>
          </div>
        ) : (
          <Outlet />
        )}
      </div>
    </div>
  );
}

export default AdminDashboard;
