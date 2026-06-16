import React from 'react';
import { Link } from 'react-router-dom';
import { Product } from '../services/api';

interface ProductCardProps {
  product: Product;
}

function ProductCard({ product }: ProductCardProps) {
  return (
    <div className="product-card">
      <Link to={`/products/${product.id}`}>
        <div className="product-card-image">
          {product.image_path ? (
            <img src={product.image_path} alt={product.name} />
          ) : (
            <div className="product-card-placeholder">No Image</div>
          )}
        </div>
        <div className="product-card-body">
          <h3>{product.name}</h3>
          <p className="product-card-price">${product.price.toFixed(2)}</p>
          <p className="product-card-stock">
            {product.stock > 0 ? `${product.stock} in stock` : 'Out of stock'}
          </p>
        </div>
      </Link>
    </div>
  );
}

export default ProductCard;
