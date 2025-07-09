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
  const response = await fetch(`${API_BASE_URL}/test-suites?project_id=${projectId}`);
  if (!response.ok) {
    throw new Error("Failed to fetch suites");
  }
  return response.json();
};

export const fetchBuilds = async (
  projectId: string | number,
  suiteId?: string | number,
): Promise<Build[]> => {
  let url = `${API_BASE_URL}/builds?project_id=${projectId}`;
  if (suiteId) {
    url += `&suite_id=${suiteId}`;
  }
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error("Failed to fetch builds");
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

export const getBuildDurationTrends = async (
  projectId: number,
  suiteId: number,
): Promise<BuildDurationTrend[]> => {
  const builds = await fetchBuilds(projectId, suiteId);
  return builds.map(build => ({
    build_number: build.build_number,
    duration: build.duration || 0,
    created_at: build.created_at,
  }));
};

export const fetchMostFailedTests = async (
  projectId: number,
  limit: number,
  suiteId?: number,
): Promise<MostFailedTest[]> => {
  const builds = await fetchBuilds(projectId, suiteId);
  const executions = await Promise.all(
    builds.map(build => fetchExecutions(build.id)),
  );

  const failedExecutions = executions
    .flat()
    .filter(execution => execution.status === 'failed');

  const failureCounts = failedExecutions.reduce(
    (acc, execution) => {
      acc[execution.test_case_id] = (acc[execution.test_case_id] || 0) + 1;
      return acc;
    },
    {} as Record<number, number>,
  );

  const testCases = await Promise.all(
    Object.keys(failureCounts).map(testCaseId =>
      fetch(`${API_BASE_URL}/test-cases?id=${testCaseId}`).then(res => res.json()),
    ),
  );

  const mostFailedTests = testCases.map(testCase => ({
    test_case_id: testCase.id,
    name: testCase.name,
    classname: testCase.classname,
    failure_count: failureCounts[testCase.id],
  }));

  return mostFailedTests
    .sort((a, b) => b.failure_count - a.failure_count)
    .slice(0, limit);
};
