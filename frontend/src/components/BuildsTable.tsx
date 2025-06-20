import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Table, Spinner, Alert, Badge } from 'react-bootstrap';
import type { Build } from '../types';
import { fetchBuilds } from '../services/api';

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
  navigate(`/projects/${projectId}/suites/${suiteId}/builds/${buildId}`);
  };

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ height: '80vh' }}>
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading builds...</span>
        </Spinner>
      </div>
    );
  }

  if (error) {
    return <Alert variant="danger">Error: {error}</Alert>;
  }

  return (
    <div className="py-3">
      <h2 className="mb-3">Builds</h2>
      <Table striped bordered hover responsive>
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
              <td className="font-monospace">{build.build_number}</td>
              <td>
                <Badge bg={getCIProviderBadgeColor(build.ci_provider)}>{build.ci_provider || 'N/A'}</Badge>
              </td>
              <td className="text-muted fst-italic">{new Date(build.created_at).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </Table>
      {builds.length === 0 && !loading && (
        <Alert variant="info" className="mt-3">No builds found for this project.</Alert>
      )}
    </div>
  );
};

// Helper function to determine badge color based on CI provider
const getCIProviderBadgeColor = (provider: string | null | undefined) => {
  if (!provider) return 'secondary';
  const lowerProvider = provider.toLowerCase();
  if (lowerProvider.includes('github')) return 'dark';
  if (lowerProvider.includes('jenkins')) return 'danger';
  if (lowerProvider.includes('travis')) return 'info';
  return 'secondary';
};

export default BuildsTable;
