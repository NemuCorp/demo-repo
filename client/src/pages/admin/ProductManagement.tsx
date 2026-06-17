import { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import * as api from '../../services/api';
import { Product } from '../../types';

function ProductManagement() {
  const { id: editId } = useParams<{ id?: string }>();
  const isEditing = !!editId;

  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [imagePath, setImagePath] = useState('');
  const [stock, setStock] = useState('');
  const [submitting, setSubmitting] = useState(false);

  const fetchProducts = () => {
    api.getProducts()
      .then((data) => setProducts(data.products))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  useEffect(() => {
    if (isEditing && editId) {
      api.getProduct(parseInt(editId))
        .then((data) => {
          const p = data.product;
          setName(p.name);
          setDescription(p.description || '');
          setPrice(p.price.toString());
          setImagePath(p.image_path || '');
          setStock(p.stock.toString());
        })
        .catch((err) => setError(err.message));
    }
  }, [editId, isEditing]);

  const resetForm = () => {
    setName('');
    setDescription('');
    setPrice('');
    setImagePath('');
    setStock('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setSuccess('');
    setSubmitting(true);

    try {
      await api.createProduct({
        name,
        description,
        price: parseFloat(price),
        image_path: imagePath,
        stock: parseInt(stock) || 0,
      });
      setSuccess(isEditing ? 'Product updated!' : 'Product created!');
      resetForm();
      fetchProducts();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return <div className="page"><p>Loading products...</p></div>;

  return (
    <div className="page admin-products-page">
      <h1>{isEditing ? 'Edit Product' : 'Product Management'}</h1>
      <Link to="/admin" className="btn btn-secondary">&larr; Back to Dashboard</Link>

      {error && <p className="error">{error}</p>}
      {success && <p className="success">{success}</p>}

      <div className="admin-form-section">
        <h2>{isEditing ? 'Edit Product' : 'Add New Product'}</h2>
        <form onSubmit={handleSubmit} className="admin-form">
          <div className="form-group">
            <label htmlFor="name">Name</label>
            <input
              id="name"
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label htmlFor="description">Description</label>
            <textarea
              id="description"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
            />
          </div>
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="price">Price</label>
              <input
                id="price"
                type="number"
                step="0.01"
                min="0"
                value={price}
                onChange={(e) => setPrice(e.target.value)}
                required
              />
            </div>
            <div className="form-group">
              <label htmlFor="stock">Stock</label>
              <input
                id="stock"
                type="number"
                min="0"
                value={stock}
                onChange={(e) => setStock(e.target.value)}
              />
            </div>
          </div>
          <div className="form-group">
            <label htmlFor="imagePath">Image URL</label>
            <input
              id="imagePath"
              type="text"
              value={imagePath}
              onChange={(e) => setImagePath(e.target.value)}
              placeholder="https://example.com/image.jpg"
            />
          </div>
          <button className="btn btn-primary" type="submit" disabled={submitting}>
            {submitting ? 'Saving...' : isEditing ? 'Update Product' : 'Create Product'}
          </button>
        </form>
      </div>

      <h2>All Products</h2>
      {products.length === 0 ? (
        <p>No products yet. Create one above.</p>
      ) : (
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Price</th>
              <th>Stock</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {products.map((p) => (
              <tr key={p.id}>
                <td>{p.id}</td>
                <td>{p.name}</td>
                <td>${p.price.toFixed(2)}</td>
                <td>{p.stock}</td>
                <td>
                  <Link to={`/admin/products/${p.id}`}>Edit</Link>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default ProductManagement;
