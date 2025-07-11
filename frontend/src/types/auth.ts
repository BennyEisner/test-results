export interface User {
  id: number;
  provider: string;
  provider_id: string;
  email: string;
  name: string;
  first_name: string;
  last_name: string;
  avatar_url: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface APIKey {
  id: number;
  user_id: number;
  name: string;
  last_used_at: string;
  expires_at: string;
  created_at: string;
  updated_at: string;
}

export interface AuthSession {
  id: string;
  user_id: number;
  provider: string;
  expires_at: string;
  created_at: string;
}

export interface AuthContext {
  user_id: number;
  provider: string;
  is_api_key: boolean;
  api_key_id?: number;
}

export interface CreateAPIKeyRequest {
  name: string;
}

export interface CreateAPIKeyResponse {
  api_key: APIKey;
  plain_text_key: string;
} 