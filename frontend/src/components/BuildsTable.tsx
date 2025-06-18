import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom'; // Added import
import type { Build } from '../types';
import { fetchBuilds } from '../services/api';
import './BuildsTable.css';

interface BuildsTableProps {
  projectId: string | number;
  suiteId: string | number;
}

const BuildsTable = ({ projectId, suiteId }: BuildsTableProps) => {
  const [builds, setBuilds] = useState<Build[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate(); // Initialize navigate

  useEffect(() => {
    const loadBuilds = async () => {
      try {
        setLoading(true);
        // Pass projectId and suiteId to fetchBuilds
        const data = await fetchBuilds(projectId, suiteId);
        setBuilds(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load builds');
      } finally {
        setLoading(false);
      }
    };

    loadBuilds();
  }, [projectId, suiteId]); 

  const handleBuildClick = (buildId: string | number) => {
    // Navigate to the executions page for the clicked build
    navigate(`/builds/${buildId}/executions`); 
  };

  if (loading) {
    return <div className="loading">Loading builds...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <h2>Builds</h2>
      <table className="table table-striped table-bordered table-hover">
        <thead>
          <tr>
            <th>Build ID</th>
            <th>Build Number</th>
            <th>CI Provider</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {builds.map((build) => (
            <tr key={build.id} onClick={() => handleBuildClick(build.id)} style={{ cursor: 'pointer' }}>
              <td>#{build.id}</td>
              <td>{build.build_number}</td>
              <td>{build.ci_provider}</td>
              <td>{new Date(build.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {builds.length === 0 && (
        <p className="no-data">No builds found for this project.</p>
      )}
    </div>
  );
};

export default BuildsTable;
