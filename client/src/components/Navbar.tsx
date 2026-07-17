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
            <Link to="/cart">Cart</Link>
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
