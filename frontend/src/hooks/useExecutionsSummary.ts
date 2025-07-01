import { useState, useEffect, useMemo } from "react";
import { fetchExecutions } from "../services/api";
import type { TestCaseExecution } from "../types";

export const useExecutionsSummary = (buildId?: string | number) => {
  const [executions, setExecutions] = useState<TestCaseExecution[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!buildId) {
      setExecutions([]);
      setLoading(false);
      return;
    }

    const loadExecutions = async () => {
      try {
        setLoading(true);
        const data = await fetchExecutions(buildId);
        setExecutions(data);
        setError(null);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to fetch executions",
        );
      } finally {
        setLoading(false);
      }
    };

    loadExecutions();
  }, [buildId]);

  const stats = useMemo(() => {
    if (loading || executions.length === 0) {
      return {
        total: 0,
        passed: 0,
        failed: 0,
        skipped: 0,
        passRate: 0,
        avgTime: 0,
      };
    }
    // Computational logic to determine build test case split
    const total = executions.length;
    const passed = executions.filter(
      (e) => e.status?.toUpperCase() === "PASSED" && e.failure == null,
    ).length;
    const failed = executions.filter(
      (e) => e.status?.toUpperCase() === "FAILED" || e.failure != null,
    ).length;
    const skipped = executions.filter(
      (e) => e.status?.toUpperCase() === "SKIPPED",
    ).length;
    const passRate = total > 0 ? Math.round((passed / total) * 100) : 0;
    const totalTime = executions.reduce(
      (sum, e) => sum + (e.execution_time || 0),
      0,
    );
    const avgTime = total > 0 ? totalTime / total : 0;

    return { total, passed, failed, skipped, passRate, avgTime };
  }, [executions, loading]);

  return { stats, executions, loading, error };
};
