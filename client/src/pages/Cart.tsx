import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import * as api from '../services/api';
import { CartItem } from '../types';

function Cart() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const fetchCart = () => {
    setLoading(true);
    api.getCart()
      .then((data) => setItems(data.cart))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchCart();
  }, []);

  const handleUpdateQuantity = async (productId: number, quantity: number) => {
    try {
      await api.updateCartItem(productId, quantity);
      fetchCart();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const handleRemove = async (productId: number) => {
    try {
      await api.removeCartItem(productId);
      fetchCart();
    } catch (err: any) {
      setError(err.message);
    }
  };

  const total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);

  if (loading) return <div className="page"><p>Loading cart...</p></div>;

  return (
    <div className="page cart-page">
      <h1>Shopping Cart</h1>
      {error && <p className="error">{error}</p>}
      {items.length === 0 ? (
        <div className="empty-cart">
          <p>Your cart is empty.</p>
          <Link to="/products" className="btn btn-primary">Browse Products</Link>
        </div>
      ) : (
        <>
          <div className="cart-items">
            {items.map((item) => (
              <div key={item.id} className="cart-item">
                <div className="cart-item-info">
                  <h3>
                    <Link to={`/products/${item.product_id}`}>{item.product_name}</Link>
                  </h3>
                  <p>${item.price.toFixed(2)} each</p>
                </div>
                <div className="cart-item-actions">
                  <label>
                    Qty:
                    <input
                      type="number"
                      min="0"
                      value={item.quantity}
                      onChange={(e) => {
                        const qty = parseInt(e.target.value) || 0;
                        if (qty === 0) {
                          handleRemove(item.product_id);
                        } else {
                          handleUpdateQuantity(item.product_id, qty);
                        }
                      }}
                    />
                  </label>
                  <p className="cart-item-subtotal">
                    ${(item.price * item.quantity).toFixed(2)}
                  </p>
                  <button
                    className="btn btn-danger"
                    onClick={() => handleRemove(item.product_id)}
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
          <div className="cart-total">
            <h2>Total: ${total.toFixed(2)}</h2>
          </div>
        </>
      )}
    </div>
  );
}

export default Cart;
