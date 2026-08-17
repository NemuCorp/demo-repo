import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { getAdminSiteUrl, getNormalSiteUrl } from '../utils/host';

function Navbar() {
  const { isAuthenticated, isAdmin, isAdminView, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate('/');
  };

  if (isAdminView) {
    return (
      <nav className="navbar admin-navbar">
        <div className="navbar-brand">
          <Link to="/admin">Demo Store Admin</Link>
        </div>
        <div className="navbar-links">
          <Link to="/admin">Dashboard</Link>
          <Link to="/admin/products">Manage Products</Link>
          <a className="btn btn-primary" href={getNormalSiteUrl()}>View Site</a>
          {isAuthenticated ? (
            <button className="btn-link" onClick={handleLogout}>Logout</button>
          ) : (
            <Link to="/login">Login</Link>
          )}
        </div>
      </nav>
    );
  }

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <Link to="/">Demo Store</Link>
      </div>
      <div className="navbar-links">
        <Link to="/products">Products</Link>
        {isAuthenticated ? (
          <>
            <Link to="/cart" className="cart-icon" title="Cart" aria-label="Cart">
              <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                <path d="M7 18c-1.1 0-1.99.9-1.99 2S5.9 22 7 22s2-.9 2-2-.9-2-2-2zM1 2v2h2l3.6 7.59-1.35 2.45c-.16.28-.25.61-.25.96 0 1.1.9 2 2 2h12v-2H7.42c-.14 0-.25-.11-.25-.25l.03-.12.9-1.63h7.45c.75 0 1.41-.41 1.75-1.03l3.58-6.49c.08-.14.12-.31.12-.48 0-.55-.45-1-1-1H5.21l-.94-2H1zm16 16c-1.1 0-1.99.9-1.99 2s.89 2 1.99 2 2-.9 2-2-.9-2-2-2z" />
              </svg>
              <span className="cart-icon-label">Cart</span>
            </Link>
            {isAdmin && (
              <a href={`${getAdminSiteUrl()}/admin`}>Admin</a>
            )}
            <button className="btn-link" onClick={handleLogout}>Logout</button>
          </>
        ) : (
          <>
            <Link to="/login">Login</Link>
            <Link to="/register">Register</Link>
          </>
        )}
      </div>
    </nav>
  );
}

export default Navbar;
