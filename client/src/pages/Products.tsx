import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import * as api from '../services/api';
import { Product } from '../types';
import { trackCartAdd, trackPageView } from '../services/tracking';

function Products() {
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [addingIds, setAddingIds] = useState<Record<number, boolean>>({});
  const [addError, setAddError] = useState('');

  useEffect(() => {
    trackPageView('/products');
    api.getProducts()
      .then((data) => setProducts(data.products))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  const handleAddToCart = async (product: Product) => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    if (product.stock === 0) return;
    setAddingIds((prev) => ({ ...prev, [product.id]: true }));
    try {
      await api.addToCart(product.id, 1);
      trackCartAdd(product.id, product.name, 1);
      setAddError('');
    } catch (err: any) {
      setAddError(err.message);
    } finally {
      setAddingIds((prev) => ({ ...prev, [product.id]: false }));
    }
  };

  if (loading) return <div className="page"><p>Loading products...</p></div>;
  if (error) return <div className="page"><p className="error">{error}</p></div>;

  return (
    <div className="page products-page">
      <h1>Products</h1>
      {products.length === 0 ? (
        <p>No products available yet.</p>
      ) : (
        <div className="product-grid">
          {products.map((product) => (
            <div key={product.id} className="product-card">
              {product.image_path && (
                <div className="product-image">
                  <img src={product.image_path} alt={product.name} />
                </div>
              )}
              <div className="product-info">
                <h3>
                  <Link to={`/products/${product.id}`}>{product.name}</Link>
                </h3>
                <p className="product-price">${product.price.toFixed(2)}</p>
                <p className="product-stock">
                  {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
                </p>
                {addError && <p className="error">{addError}</p>}
                <button
                  className="btn btn-primary btn-sm btn-full add-to-cart-btn"
                  disabled={product.stock === 0 || addingIds[product.id]}
                  onClick={() => handleAddToCart(product)}
                >
                  {addingIds[product.id] ? 'Adding...' : 'Add to Cart'}
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

export default Products;
