import type { Project } from "../types";

import type { Build } from "../types";

const API_BASE_URL = "http://localhost:8080/api";

export const fetchProjects = async (): Promise<Project[]> => {
  const response = await fetch(`${API_BASE_URL}/projects`);
  if (!response.ok) {
    throw new Error(`Failed to fetch projects: ${response.status}`);
  }
  return response.json();
};

export const fetchBuilds = async (
  projectId: string | number,
): Promise<Build[]> => {
  const response = await fetch(`/api/projects/${projectId}/builds`);
  if (!response.ok) {
    throw new Error("Failed to fetch builds");
  }
  return response.json();
};
