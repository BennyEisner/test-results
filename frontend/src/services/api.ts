import type { Project } from "../types";

import type { Suite } from "../types";

import type { Build } from "../types";

import type { TestCaseExecution } from "../types";

import type { Failure, SearchResult, BuildDurationTrend, MostFailedTest } from "../types";

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
  suiteId?: string | number,
): Promise<Build[]> => {
  let url = `${API_BASE_URL}/projects/${projectId}/builds`;
  if (suiteId) {
    url = `${API_BASE_URL}/projects/${projectId}/suites/${suiteId}/builds`;
  }
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Failed to fetch builds");
  }
  return response.json();
};

export const fetchRecentBuilds = async (projectId?: string | number): Promise<Build[]> => {
  let url = `${API_BASE_URL}/builds/recent`;
  if (projectId) {
    url = `${API_BASE_URL}/projects/${projectId}/builds/recent`;
  }
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Failed to fetch recent builds");
  }
  return response.json();
};

export const fetchExecutions = async (
  buildId: string | number,
): Promise<TestCaseExecution[]> => {
  // Used TestCaseExecution type
  const response = await fetch(`${API_BASE_URL}/builds/${buildId}/executions`); // Added API_BASE_URL
  if (!response.ok) {
    throw new Error("Failed to fetch executions");
  }
  return response.json();
};

export const fetchFailures = async (
  buildId: string | number,
): Promise<Failure[]> => {
  const response = await fetch(`${API_BASE_URL}/builds/${buildId}/failures`); // Added API_BASE_URL
  if (!response.ok) {
    throw new Error("Failed to fetch executions");
  }
  return response.json();
};

export const search = async (query: string): Promise<SearchResult[]> => {
  if (!query.trim()) {
    return [];
  }
  const response = await fetch(`${API_BASE_URL}/search?q=${encodeURIComponent(query)}`);
  if (!response.ok) {
    throw new Error("Failed to search");
  }
  return response.json();
};

export const getBuildDurationTrends = async (projectId: number, suiteId: number): Promise<BuildDurationTrend[]> => {
  const response = await fetch(`${API_BASE_URL}/builds/duration-trends?projectId=${projectId}&suiteId=${suiteId}`);
  if (!response.ok) {
    throw new Error('Failed to fetch build duration trends');
  }
  return response.json();
};

export const fetchMostFailedTests = async (projectId: number, limit: number): Promise<MostFailedTest[]> => {
  const response = await fetch(`${API_BASE_URL}/test-cases/most-failed?projectId=${projectId}&limit=${limit}`);
  if (!response.ok) {
    throw new Error('Failed to fetch most failed tests');
  }
  return response.json();
};
