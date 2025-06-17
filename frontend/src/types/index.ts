// TypeScript interfaces matching the Go backend models
export interface Project {
  id: number;
  name: string;
}

export interface Build {
  id: number;
  project_id: number;
  build_number: string;
  ci_provider: string;
  ci_url?: string;
  created_at: string;
}

export interface Suite {
  id: number;
  project_id: number;
  name: string;
  parent_id?: number;
  time: number;
}
