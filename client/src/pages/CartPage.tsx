import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { CartItem, getCart, updateCartItem, removeFromCart } from '../services/api';

function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchCart = () => {
    setLoading(true);
    getCart()
      .then((data) => setItems(data.cart))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchCart();
  }, []);

  const handleUpdateQuantity = async (productId: number, quantity: number) => {
    try {
      await updateCartItem(productId, quantity);
      fetchCart();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update');
    }
  };

  const handleRemove = async (productId: number) => {
    try {
      await removeFromCart(productId);
      fetchCart();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to remove');
    }
  };

  const total = items.reduce((sum, item) => sum + item.price * item.quantity, 0);

  if (loading) {
    return <div className="loading">Loading cart...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="page">
      <h1>Shopping Cart</h1>
      {items.length === 0 ? (
        <p>
          Your cart is empty. <Link to="/">Browse products</Link>
        </p>
      ) : (
        <>
          <div className="cart-items">
            {items.map((item) => (
              <div key={item.id} className="cart-item">
                <div className="cart-item-info">
                  <Link to={`/products/${item.product_id}`} className="cart-item-name">
                    {item.product_name}
                  </Link>
                  <span className="cart-item-price">${item.price.toFixed(2)} each</span>
                </div>
                <div className="cart-item-actions">
                  <label>
                    Qty:
                    <input
                      type="number"
                      min={0}
                      value={item.quantity}
                      onChange={(e) => {
                        const qty = parseInt(e.target.value, 10) || 0;
                        if (qty === 0) {
                          handleRemove(item.product_id);
                        } else {
                          handleUpdateQuantity(item.product_id, qty);
                        }
                      }}
                    />
                  </label>
                  <span className="cart-item-subtotal">
                    ${(item.price * item.quantity).toFixed(2)}
                  </span>
                  <button
                    onClick={() => handleRemove(item.product_id)}
                    className="btn-danger btn-sm"
                  >
                    Remove
                  </button>
                </div>
              </div>
            ))}
          </div>
          <div className="cart-total">
            <strong>Total: ${total.toFixed(2)}</strong>
          </div>
        </>
      )}
    </div>
  );
}

export default CartPage;
