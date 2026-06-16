import React, { useEffect, useState, FormEvent } from 'react';
import { Product, getProducts, createProduct } from '../../services/api';

function AdminProducts() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [price, setPrice] = useState('');
  const [imagePath, setImagePath] = useState('');
  const [stock, setStock] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitSuccess, setSubmitSuccess] = useState(false);

  const fetchProducts = () => {
    setLoading(true);
    getProducts()
      .then((data) => setProducts(data.products))
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  };

  useEffect(() => {
    fetchProducts();
  }, []);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setSubmitError(null);
    setSubmitSuccess(false);

    try {
      await createProduct({
        name,
        description,
        price: parseFloat(price) || 0,
        image_path: imagePath,
        stock: parseInt(stock, 10) || 0,
      });
      setSubmitSuccess(true);
      setName('');
      setDescription('');
      setPrice('');
      setImagePath('');
      setStock('');
      setShowForm(false);
      fetchProducts();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : 'Failed to create product');
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="loading">Loading products...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <div className="admin-header">
        <h2>Products ({products.length})</h2>
        <button
          onClick={() => { setShowForm(!showForm); setSubmitError(null); setSubmitSuccess(false); }}
          className="btn-primary"
        >
          {showForm ? 'Cancel' : 'Add Product'}
        </button>
      </div>

      {showForm && (
        <form onSubmit={handleSubmit} className="admin-form">
          <h3>New Product</h3>
          <div className="form-group">
            <label htmlFor="pname">Name</label>
            <input id="pname" type="text" value={name} onChange={(e) => setName(e.target.value)} required />
          </div>
          <div className="form-group">
            <label htmlFor="pdesc">Description</label>
            <textarea id="pdesc" value={description} onChange={(e) => setDescription(e.target.value)} />
          </div>
          <div className="form-row">
            <div className="form-group">
              <label htmlFor="pprice">Price</label>
              <input id="pprice" type="number" step="0.01" min="0" value={price} onChange={(e) => setPrice(e.target.value)} required />
            </div>
            <div className="form-group">
              <label htmlFor="pstock">Stock</label>
              <input id="pstock" type="number" min="0" value={stock} onChange={(e) => setStock(e.target.value)} />
            </div>
          </div>
          <div className="form-group">
            <label htmlFor="pimage">Image Path</label>
            <input id="pimage" type="text" value={imagePath} onChange={(e) => setImagePath(e.target.value)} />
          </div>
          {submitError && <div className="error-msg">{submitError}</div>}
          {submitSuccess && <div className="success-msg">Product created successfully!</div>}
          <button type="submit" disabled={submitting} className="btn-primary">
            {submitting ? 'Creating...' : 'Create Product'}
          </button>
        </form>
      )}

      {products.length === 0 ? (
        <p>No products yet.</p>
      ) : (
        <table className="admin-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>Name</th>
              <th>Price</th>
              <th>Stock</th>
            </tr>
          </thead>
          <tbody>
            {products.map((p) => (
              <tr key={p.id}>
                <td>{p.id}</td>
                <td>{p.name}</td>
                <td>${p.price.toFixed(2)}</td>
                <td>{p.stock}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default AdminProducts;
