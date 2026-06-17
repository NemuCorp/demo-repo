import React, { useEffect } from 'react';
import { Link, Outlet, useLocation } from 'react-router-dom';
import AdminStats from './AdminStats';
import { trackPageView } from '../../services/tracking';

function AdminDashboard() {
  const location = useLocation();

  useEffect(() => {
    trackPageView(location.pathname);
  }, [location.pathname]);

  const isRoot = location.pathname === '/admin';
  const isProducts = location.pathname === '/admin/products';
  const isActivity = location.pathname === '/admin/activity';

  return (
    <div className="page">
      <h1>Admin</h1>
      <div className="admin-nav">
        <Link
          to="/admin"
          className={isRoot ? 'active' : ''}
        >
          Dashboard
        </Link>
        <Link
          to="/admin/products"
          className={isProducts ? 'active' : ''}
        >
          Manage Products
        </Link>
        <Link
          to="/admin/activity"
          className={isActivity ? 'active' : ''}
        >
          Activity
        </Link>
      </div>
      <div className="admin-content">
        {isRoot ? (
          <AdminStats detailed={false} />
        ) : isActivity ? (
          <AdminStats detailed={true} />
        ) : (
          <Outlet />
        )}
      </div>
    </div>
  );
}

export default AdminDashboard;
