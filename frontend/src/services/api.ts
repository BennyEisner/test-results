import axios from 'axios';
import type { Project, Suite, Build, TestCaseExecution, Failure, SearchResult, BuildDurationTrend, MostFailedTest } from "../types";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  withCredentials: true,
});

export const fetchProjects = async (): Promise<Project[]> => {
  const response = await api.get(`/projects`);
  return response.data;
};

export const fetchSuites = async (
  projectId: string | number,
): Promise<Suite[]> => {
  const response = await api.get(`/test-suites?project_id=${projectId}`);
  return response.data;
};

export const fetchBuilds = async (
  projectId: string | number,
  suiteId?: string | number,
): Promise<Build[]> => {
  let url = `/builds?project_id=${projectId}`;
  if (suiteId) {
    url += `&suite_id=${suiteId}`;
  }
  const response = await api.get(url);
  return response.data;
};

export const fetchExecutions = async (
  buildId: string | number,
): Promise<TestCaseExecution[]> => {
  const response = await api.get(`/builds/${buildId}/executions`);
  return response.data;
};

export const fetchFailures = async (
  buildId: string | number,
): Promise<Failure[]> => {
  const response = await api.get(`/builds/${buildId}/failures`);
  return response.data;
};

export const search = async (query: string): Promise<SearchResult[]> => {
  if (!query.trim()) {
    return [];
  }
  const response = await api.get(`/search?q=${encodeURIComponent(query)}`);
  return response.data;
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
      api.get(`/test-cases?id=${testCaseId}`).then(res => res.data),
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

export default api;
