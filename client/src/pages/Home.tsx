import { Link } from 'react-router-dom';

function Home() {
  return (
    <div className="page home-page">
      <section className="hero">
        <h1>Welcome to Demo Store</h1>
        <p>Your one-stop shop for everything you need.</p>
        <Link to="/products" className="btn btn-primary">Shop Now</Link>
      </section>
      <section className="features">
        <div className="feature-card">
          <h3>Wide Selection</h3>
          <p>Browse our catalog of quality products at competitive prices.</p>
        </div>
        <div className="feature-card">
          <h3>Easy Cart</h3>
          <p>Add items to your cart and manage quantities with ease.</p>
        </div>
        <div className="feature-card">
          <h3>Secure Account</h3>
          <p>Create an account to track orders and save your preferences.</p>
        </div>
      </section>
    </div>
  );
}

export default Home;
