const TRACK_ENABLED = true;

function getToken(): string | null {
  return localStorage.getItem('auth_token');
}

async function sendEvent(payload: { event_type: string; event_data?: Record<string, unknown> }): Promise<void> {
  if (!TRACK_ENABLED) return;
  try {
    const token = getToken();
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    };
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }
    await fetch('/api/track', {
      method: 'POST',
      headers,
      body: JSON.stringify(payload),
    });
  } catch {
    /* silently ignore tracking failures */
  }
}

export function trackPageView(path: string): void {
  sendEvent({
    event_type: 'page_view',
    event_data: { path, title: document.title },
  });
}

export function trackProductView(productId: number, productName: string): void {
  sendEvent({
    event_type: 'product_view',
    event_data: { product_id: String(productId), product_name: productName },
  });
}

export function trackCartAdd(productId: number, productName: string, quantity: number): void {
  sendEvent({
    event_type: 'cart_add',
    event_data: { product_id: String(productId), product_name: productName, quantity: String(quantity) },
  });
}

export function trackCartRemove(productId: number, productName: string): void {
  sendEvent({
    event_type: 'cart_remove',
    event_data: { product_id: String(productId), product_name: productName },
  });
}

export function trackRegistration(): void {
  sendEvent({ event_type: 'registration', event_data: {} });
}

export function trackLogin(): void {
  sendEvent({ event_type: 'login', event_data: {} });
}

export function trackProductCreate(productName: string): void {
  sendEvent({
    event_type: 'product_created',
    event_data: { product_name: productName },
  });
}
