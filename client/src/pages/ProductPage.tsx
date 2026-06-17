import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Product, getProduct, addToCart } from '../services/api';
import { trackPageView, trackProductView, trackCartAdd } from '../services/tracking';
import { useAuth } from '../contexts/AuthContext';

function ProductPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [product, setProduct] = useState<Product | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [adding, setAdding] = useState(false);
  const [addError, setAddError] = useState<string | null>(null);
  const [addSuccess, setAddSuccess] = useState(false);

  useEffect(() => {
    if (!id) return;
    const productId = parseInt(id, 10);
    getProduct(productId)
      .then((data) => {
        setProduct(data.product);
        trackProductView(data.product.id, data.product.name);
        trackPageView(`/products/${id}`);
      })
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [id]);

  const handleAddToCart = async () => {
    if (!user) {
      navigate('/login');
      return;
    }
    if (!product) return;

    setAdding(true);
    setAddError(null);
    setAddSuccess(false);
    try {
      await addToCart(product.id, quantity);
      trackCartAdd(product.id, product.name, quantity);
      setAddSuccess(true);
    } catch (err) {
      setAddError(err instanceof Error ? err.message : 'Failed to add to cart');
    } finally {
      setAdding(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading product...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  if (!product) {
    return <div className="error">Product not found.</div>;
  }

  return (
    <div className="page">
      <button onClick={() => navigate(-1)} className="btn-link">&larr; Back</button>
      <div className="product-detail">
        <div className="product-detail-image">
          {product.image_path ? (
            <img src={product.image_path} alt={product.name} />
          ) : (
            <div className="product-card-placeholder large">No Image</div>
          )}
        </div>
        <div className="product-detail-info">
          <h1>{product.name}</h1>
          <p className="product-detail-price">${product.price.toFixed(2)}</p>
          <p className="product-detail-stock">
            {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
          </p>
          {product.description && <p className="product-detail-desc">{product.description}</p>}

          {product.stock > 0 && (
            <div className="product-detail-actions">
              <label>
                Quantity:
                <input
                  type="number"
                  min={1}
                  max={product.stock}
                  value={quantity}
                  onChange={(e) => {
                    const v = parseInt(e.target.value, 10);
                    setQuantity(isNaN(v) || v < 1 ? 1 : v > product.stock ? product.stock : v);
                  }}
                />
              </label>
              <button
                onClick={handleAddToCart}
                disabled={adding}
                className="btn-primary"
              >
                {adding ? 'Adding...' : 'Add to Cart'}
              </button>
              {addSuccess && <span className="success-msg">Added to cart!</span>}
              {addError && <span className="error-msg">{addError}</span>}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

export default ProductPage;
