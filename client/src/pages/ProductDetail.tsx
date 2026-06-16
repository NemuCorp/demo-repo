import { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import * as api from '../services/api';
import { Product } from '../types';

function ProductDetail() {
  const { id } = useParams<{ id: string }>();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const [product, setProduct] = useState<Product | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [quantity, setQuantity] = useState(1);
  const [adding, setAdding] = useState(false);

  useEffect(() => {
    if (!id) return;
    api.getProduct(parseInt(id))
      .then((data) => setProduct(data.product))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [id]);

  const handleAddToCart = async () => {
    if (!isAuthenticated) {
      navigate('/login');
      return;
    }
    if (!product) return;
    setAdding(true);
    try {
      await api.addToCart(product.id, quantity);
      alert('Added to cart!');
    } catch (err: any) {
      setError(err.message);
    } finally {
      setAdding(false);
    }
  };

  if (loading) return <div className="page"><p>Loading product...</p></div>;
  if (error) return <div className="page"><p className="error">{error}</p></div>;
  if (!product) return <div className="page"><p>Product not found.</p></div>;

  return (
    <div className="page product-detail-page">
      <div className="product-detail">
        {product.image_path && (
          <div className="product-detail-image">
            <img src={product.image_path} alt={product.name} />
          </div>
        )}
        <div className="product-detail-info">
          <h1>{product.name}</h1>
          <p className="product-price">${product.price.toFixed(2)}</p>
          <p className="product-description">{product.description}</p>
          <p className={`product-stock ${product.stock === 0 ? 'out-of-stock' : ''}`}>
            {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
          </p>
          {product.stock > 0 && (
            <div className="add-to-cart-form">
              <label>
                Quantity:
                <input
                  type="number"
                  min="1"
                  max={product.stock}
                  value={quantity}
                  onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                />
              </label>
              <button
                className="btn btn-primary"
                onClick={handleAddToCart}
                disabled={adding}
              >
                {adding ? 'Adding...' : 'Add to Cart'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ProductDetail;
