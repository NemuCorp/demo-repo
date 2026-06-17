import { AuthResponse, CartItem, Product, User } from '../types';

export type { CartItem, Product };

const API_BASE = process.env.REACT_APP_API_URL || '/api';

function getToken(): string | null {
  return localStorage.getItem('auth_token');
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(data.error || 'An error occurred');
  }

  return data as T;
}

export async function register(email: string, password: string): Promise<{ user: User }> {
  return request<{ user: User }>('/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export async function login(email: string, password: string): Promise<AuthResponse> {
  return request<AuthResponse>('/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export async function logout(): Promise<{ message: string }> {
  return request<{ message: string }>('/auth/logout', {
    method: 'POST',
  });
}

export async function getProducts(): Promise<{ products: Product[] }> {
  return request<{ products: Product[] }>('/products');
}

export async function getProduct(id: number): Promise<{ product: Product }> {
  return request<{ product: Product }>(`/products/${id}`);
}

export async function createProduct(data: {
  name: string;
  description: string;
  price: number;
  image_path: string;
  stock: number;
}): Promise<{ product: Product }> {
  return request<{ product: Product }>('/products', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getCart(): Promise<{ cart: CartItem[] }> {
  return request<{ cart: CartItem[] }>('/cart');
}

export async function addToCart(productId: number, quantity: number): Promise<{ item: CartItem }> {
  return request<{ item: CartItem }>('/cart', {
    method: 'POST',
    body: JSON.stringify({ product_id: productId, quantity }),
  });
}

export async function updateCartItem(productId: number, quantity: number): Promise<{ item: CartItem } | { message: string }> {
  return request<{ item: CartItem } | { message: string }>(`/cart/${productId}`, {
    method: 'PUT',
    body: JSON.stringify({ quantity }),
  });
}

export async function removeCartItem(productId: number): Promise<{ message: string }> {
  return request<{ message: string }>(`/cart/${productId}`, {
    method: 'DELETE',
  });
}

export interface DashboardMetrics {
  total_users: number;
  total_products: number;
  page_views: number;
  product_views: number;
  cart_adds: number;
  registrations: number;
  today: {
    active_users: number;
    total_events: number;
    products_viewed: number;
  };
  top_products: Array<{
    product_id: string;
    product_name: string;
    views: number;
  }>;
  recent_activity: Array<{
    id: number;
    user_id: number | null;
    user_email: string;
    event_type: string;
    event_data: unknown;
    created_at: string;
  }>;
  daily_stats: Array<{
    day: string;
    event_type: string;
    count: number;
  }>;
}

export async function getDashboardMetrics(): Promise<DashboardMetrics> {
  return request<DashboardMetrics>('/api/admin/stats');
}
