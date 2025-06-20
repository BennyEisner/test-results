import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import ExecutionsTable from './ExecutionsTable';
import ExecutionsSummary from './ExecutionsSummary';
import { fetchExecutions } from '../services/api';
import type { TestCaseExecution } from '../types';
import './ExecutionsSummary.css';

const BuildDetail = () => {
  const { buildId, suiteId, projectId } = useParams<{ buildId: string, suiteId: string, projectId: string }>();
  const navigate = useNavigate();
  
  // State for execution data - lifted up from ExecutionsTable
  const [executions, setExecutions] = useState<TestCaseExecution[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Validation checks
  if (!suiteId) {
    return <div className="error">Suite ID is required</div>;
  }

  if (!projectId) {
    return <div className="error">Project ID is required</div>;
  }

  if (!buildId) {
    return <div className="error">Build ID is required</div>; 
  }

  // Fetch execution data
  useEffect(() => {
    const loadExecutions = async () => {
      try {
        setLoading(true);
        const data = await fetchExecutions(buildId);
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

  return (
    <div className="build-detail">
      <div className="build-header">
        <button 
          onClick={() => navigate(`/projects/${projectId}/suites/${suiteId}`)} 
          className="back-button"
        >
          Back to Builds
        </button>
        <h1>Build {buildId}</h1>
      </div>

      {/* New summary component */}
      <ExecutionsSummary
        executions={executions}
        loading={loading}
      />

      {/* Pass data to ExecutionsTable */}
      {error ? (
        <div className="error">Error: {error}</div>
      ) : (
        <ExecutionsTable
          executions={executions}
          loading={loading}
        />
      )}
    </div>
  );
};

export default BuildDetail;
