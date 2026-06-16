export interface User {
  id: number;
  email: string;
  created_at: string;
}

export interface Session {
  id: string;
  user_id: number;
  created_at: string;
  expires_at: string;
}

export interface Product {
  id: number;
  name: string;
  description: string;
  price: number;
  image_path: string;
  stock: number;
  created_at: string;
  updated_at: string;
}

export interface CartItem {
  id: number;
  user_id: number;
  product_id: number;
  quantity: number;
  created_at: string;
  product_name: string;
  price: number;
  image_path: string;
}

export interface ApiError {
  error: string;
}

const TOKEN_KEY = 'auth_token';

function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export function isLoggedIn(): boolean {
  return getToken() !== null;
}

async function request<T>(
  path: string,
  options: RequestInit = {}
): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...((options.headers as Record<string, string>) || {}),
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  const res = await fetch(path, {
    ...options,
    headers,
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error((data as ApiError).error || `Request failed with status ${res.status}`);
  }

  return data as T;
}

export async function register(email: string, password: string): Promise<{ user: User }> {
  return request<{ user: User }>('/api/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
}

export async function login(email: string, password: string): Promise<{ token: string; session: Session }> {
  const data = await request<{ token: string; session: Session }>('/api/auth/login', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
  setToken(data.token);
  return data;
}

export async function logout(): Promise<void> {
  try {
    await request<{ message: string }>('/api/auth/logout', { method: 'POST' });
  } finally {
    clearToken();
  }
}

export async function getProducts(): Promise<{ products: Product[] }> {
  return request<{ products: Product[] }>('/api/products');
}

export async function getProduct(id: number): Promise<{ product: Product }> {
  return request<{ product: Product }>(`/api/products/${id}`);
}

export async function createProduct(data: {
  name: string;
  description: string;
  price: number;
  image_path: string;
  stock: number;
}): Promise<{ product: Product }> {
  return request<{ product: Product }>('/api/products', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function getCart(): Promise<{ cart: CartItem[] }> {
  return request<{ cart: CartItem[] }>('/api/cart');
}

export async function addToCart(productId: number, quantity: number): Promise<{ item: CartItem }> {
  return request<{ item: CartItem }>('/api/cart', {
    method: 'POST',
    body: JSON.stringify({ product_id: productId, quantity }),
  });
}

export async function updateCartItem(productId: number, quantity: number): Promise<{ item: CartItem } | { message: string }> {
  return request<{ item: CartItem } | { message: string }>(`/api/cart/${productId}`, {
    method: 'PUT',
    body: JSON.stringify({ quantity }),
  });
}

export async function removeFromCart(productId: number): Promise<{ message: string }> {
  return request<{ message: string }>(`/api/cart/${productId}`, {
    method: 'DELETE',
  });
}
