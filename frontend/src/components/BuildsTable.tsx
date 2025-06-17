import { useState, useEffect } from 'react';
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

  useEffect(() => {
    const loadBuilds = async () => {
      try {
        setLoading(true);
        // Pass both projectId and suiteId to fetchBuilds
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
  }, [projectId, suiteId]); // Add suiteId to dependency array

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
            <th>CI URL</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {builds.map((build) => (
            <tr key={build.id}>
              <td>#{build.id}</td>
              <td>{build.build_number}</td>
              <td>{build.ci_provider}</td>
              <td>
                {build.ci_url ? (
                  <a 
                    href={build.ci_url} 
                    target="_blank" 
                    rel="noopener noreferrer"
                  >
                    View Build
                  </a>
                ) : (
                  '-'
                )}
              </td>
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
