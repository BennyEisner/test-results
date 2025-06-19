// ExecutionsSummary.tsx
import { useMemo } from 'react';
import type { TestCaseExecution } from '../types';
import './ExecutionsSummary.css';

interface ExecutionsSummaryProps {
  executions: TestCaseExecution[];
  loading: boolean;
}

const ExecutionsSummary = ({ executions, loading }: ExecutionsSummaryProps) => {
  // Calculate execution statistics only when the executions data changes
  const stats = useMemo(() => {
    if (loading || executions.length === 0) {
      return {
        total: 0,
        passed: 0,
        failed: 0,
        skipped: 0,
        passRate: 0,
        avgTime: 0
      };
    }

    // Count total executions
    const total = executions.length;
    
    // Count executions by status
    const passed = executions.filter(e => e.status === 'PASSED').length;
    const failed = executions.filter(e => 
      e.status === 'FAILED' || e.failure != null
    ).length;
    const skipped = executions.filter(e => e.status === 'SKIPPED').length;
    
    // Calculate pass rate percentage
    const passRate = total > 0 ? Math.round((passed / total) * 100) : 0;
    
    // Calculate average execution time (in seconds)
    const totalTime = executions.reduce((sum, e) => 
      sum + (e.execution_time || 0), 0
    );
    const avgTime = total > 0 ? (totalTime / total) : 0;
    
    return { total, passed, failed, skipped, passRate, avgTime };
  }, [executions, loading]);

  if (loading) {
    return <div className="summary-loading">Analyzing test results...</div>;
  }

  // Render the summary metrics
  return (
    <div className="executions-summary">
      <h3>Test Results Summary</h3>
      
      <div className="summary-grid">
        <div className="summary-card total">
          <div className="summary-value">{stats.total}</div>
          <div className="summary-label">Total Tests</div>
        </div>
        
        <div className="summary-card passed">
          <div className="summary-value">{stats.passed}</div>
          <div className="summary-label">Passed</div>
        </div>
        
        <div className="summary-card failed">
          <div className="summary-value">{stats.failed}</div>
          <div className="summary-label">Failed</div>
        </div>
        
        <div className="summary-card skipped">
          <div className="summary-value">{stats.skipped}</div>
          <div className="summary-label">Skipped</div>
        </div>
        
        <div className="summary-card rate">
          <div className="summary-value">{stats.passRate}%</div>
          <div className="summary-label">Pass Rate</div>
        </div>
        
        <div className="summary-card time">
          <div className="summary-value">{stats.avgTime.toFixed(2)}s</div>
          <div className="summary-label">Avg Time</div>
        </div>
      </div>
    </div>
  );
};

export default ExecutionsSummary;
