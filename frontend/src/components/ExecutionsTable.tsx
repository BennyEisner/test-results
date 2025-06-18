import { useState, useEffect } from 'react';
import type { TestCaseExecution } from '../types'; // Renamed Execution to TestCaseExecution
import { fetchExecutions } from '../services/api'; // Corrected import name

interface ExecutionsTableProps {
  buildId: string | number;
}

const ExecutionsTable = ({ buildId }: ExecutionsTableProps) => { 
  const [executions, setExecutions] = useState<TestCaseExecution[]>([]); // Used TestCaseExecution type
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadExecutions = async () => {
      try {
        setLoading(true);
        // Ensure buildId is passed correctly, it's an object initially
        const data = await fetchExecutions(buildId.toString()); 
        setExecutions(data); 
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load executions');
      } finally {
        setLoading(false);
      }
    };

    loadExecutions();
  }, [buildId]); 

  if (loading) {
    return <div className="loading">Loading executions...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <h2>Executions</h2>
      <table className="table table-striped table-bordered table-hover">
        <thead>
          <tr>
            <th>Execution ID</th>
            <th>Test Case Name</th>
            <th>Status</th>
            <th>Execution Time</th>
            <th>Created At</th>
          </tr>
        </thead>
        <tbody>
          {executions.map((execution) => (
            <tr key={execution.id} className={execution.failure ? 'table-danger' : ''}>
              <td>#{execution.id}</td>
              <td>{execution.test_case_name || `Test Case ${execution.test_case_id}`}</td>
              <td>
                {execution.status}
                {execution.failure && (
                  <span title={`Failure: ${execution.failure.message || 'No message'}`} style={{ marginLeft: '8px', color: 'red', cursor: 'help' }}>
                  </span>
                )}
              </td>
              <td>{execution.execution_time}</td>
              <td>{new Date(execution.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {executions.length === 0 && (
        <p className="no-data">No executions found for this project.</p>
      )}
    </div>
  );
};

export default ExecutionsTable;
