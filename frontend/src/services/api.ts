import type { Project } from "../types";

import type { Suite } from "../types";

import type { Build } from "../types";

import type { TestCaseExecution } from "../types"; // Renamed Executions to TestCaseExecution

const API_BASE_URL = "http://localhost:8080/api";

export const fetchProjects = async (): Promise<Project[]> => {
  const response = await fetch(`${API_BASE_URL}/projects`);
  if (!response.ok) {
    throw new Error(`Failed to fetch projects: ${response.status}`);
  }
  return response.json();
};

export const fetchSuites = async (
  projectId: string | number,
): Promise<Suite[]> => {
  const response = await fetch(`${API_BASE_URL}/projects/${projectId}/suites`); 
  if (!response.ok) {
    throw new Error("Failed to fetch suites");
  }
  return response.json();
};

export const fetchBuilds = async (
  projectId: string | number, 
  suiteId: string | number,
): Promise<Build[]> => {
  
  const response = await fetch(
    `${API_BASE_URL}/projects/${projectId}/suites/${suiteId}/builds`, 
  );
  if (!response.ok) {
    throw new Error("Failed to fetch builds");
  }
  return response.json();
};

export const fetchExecutions = async (
  buildId: string | number,
): Promise<TestCaseExecution[]> => { // Used TestCaseExecution type
  const response = await fetch(`${API_BASE_URL}/builds/${buildId}/executions`); // Added API_BASE_URL
  if (!response.ok) {
    throw new Error("Failed to fetch executions");
  }
  return response.json();
};
