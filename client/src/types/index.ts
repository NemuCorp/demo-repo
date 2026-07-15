export interface User {
  id: number;
  email: string;
  is_admin?: boolean;
  created_at: string;
  updated_at: string;
}

export interface Session {
  id: string;
  user_id: number;
  created_at: string;
  expires_at: string;
}

export interface AuthResponse {
  token: string;
  session: Session;
  is_admin?: boolean;
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
