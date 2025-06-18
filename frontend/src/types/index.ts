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

export interface TestCaseExecution { 
  id: number;
  build_id: number;
  test_case_id: number;
  status: string;
  execution_time: number;
  created_at: string;
  test_case_name?: string; // Added from backend's BuildExecutionDetail
  class_name?: string;   // Added from backend's BuildExecutionDetail
  failure?: Failure | null; // Added from backend's BuildExecutionDetail
}

// Added Failure interface to match backend's models.Failure (simplified for frontend)
export interface Failure {
  message?: string | null;
  type?: string | null;
  details?: string | null;
}
