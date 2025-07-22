import { User, APIKey, CreateAPIKeyRequest, CreateAPIKeyResponse } from '../types/auth';

const AUTH_BASE = '/auth';

export const authApi = {
  // OAuth2 Authentication
  beginOAuth2Auth: (provider: string): void => {
    window.location.href = `${AUTH_BASE}/${provider}`;
  },

  // Session Management
  logout: async (): Promise<void> => {
    await fetch(`${AUTH_BASE}/logout`, {
      method: 'POST',
      credentials: 'include'
    });
  },

  // User Management
  getCurrentUser: async (): Promise<User> => {
    const response = await fetch(`${AUTH_BASE}/me`, {
      credentials: 'include'
    });
    
    if (!response.ok) {
      throw new Error('Failed to get current user');
    }
    
    return response.json();
  },

  // API Key Management
  listAPIKeys: async (): Promise<APIKey[]> => {
    const response = await fetch(`${AUTH_BASE}/api-keys`, {
      credentials: 'include'
    });
    
    if (!response.ok) {
      throw new Error('Failed to list API keys');
    }
    
    return response.json();
  },

  createAPIKey: async (request: CreateAPIKeyRequest): Promise<CreateAPIKeyResponse> => {
    const response = await fetch(`${AUTH_BASE}/api-keys`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(request)
    });
    
    if (!response.ok) {
      throw new Error('Failed to create API key');
    }
    
    return response.json();
  },

  deleteAPIKey: async (keyId: number): Promise<void> => {
    const response = await fetch(`${AUTH_BASE}/api-keys/${keyId}`, {
      method: 'DELETE',
      credentials: 'include'
    });
    
    if (!response.ok) {
      throw new Error('Failed to delete API key');
    }
  }
};
