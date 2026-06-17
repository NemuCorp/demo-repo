import React, { useEffect, useState } from 'react';
import { Product, getProducts } from '../services/api';
import { trackPageView } from '../services/tracking';
import ProductCard from '../components/ProductCard';

function HomePage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    trackPageView('/');
  }, []);

  useEffect(() => {
    getProducts()
      .then((data) => setProducts(data.products))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="loading">Loading products...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div className="page">
      <h1>Products</h1>
      {products.length === 0 ? (
        <p>No products available yet.</p>
      ) : (
        <div className="product-grid">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}

export default HomePage;
