import { useState, useEffect } from 'react';
import type { Suite } from '../types';
import { fetchSuites } from '../services/api';

interface SuitesTableProps {
  projectId: string;
}

const SuitesTable = ({ projectId }: SuitesTableProps) => {
  const [suites, setSuites] = useState<Suite[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadSuites = async () => {
      try {
        setLoading(true);
        const data = await fetchSuites(projectId);
        setSuites(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load test_suites');
      } finally {
        setLoading(false);
      }
    };

    loadSuites();
  }, [projectId]);

  if (loading) {
    return <div className="loading">Loading suites...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <h2>Suites</h2>
      <table className="table table-striped table-bordered table-hover">
        <thead>
          <tr>
            <th>Suite ID</th>
            <th>Name</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {suites.map((suite) => (
            <tr key={suite.id}>
              <td>#{suite.id}</td>
              <td>{suite.name}</td>
              <td>{new Date(suite.time).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {suites.length === 0 && (
        <p className="no-data">No suites found for this project.</p>
      )}
    </div>
  );
};

export default SuitesTable;
