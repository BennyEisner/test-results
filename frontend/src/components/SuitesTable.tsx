import { useState, useEffect } from 'react';
import { Table, Spinner, Alert } from 'react-bootstrap';
import type { Suite } from '../types';
import { fetchSuites } from '../services/api';
import { useNavigate } from 'react-router-dom';

interface SuitesTableProps {
  projectId: string;
}

const SuitesTable = ({ projectId }: SuitesTableProps) => {
  const navigate = useNavigate();
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

  const handleSuiteClick = (suiteId: string | number) => {
    navigate(`/projects/${projectId}/suites/${suiteId}`);
  };

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ height: '80vh' }}>
        <Spinner animation="border" role="status">
          <span className="visually-hidden">Loading suites...</span>
        </Spinner>
      </div>
    );
  }

  if (error) {
    return <Alert variant="danger">Error: {error}</Alert>;
  }

  return (
    <div className="py-3">
      <h2 className="mb-3">Suites</h2>
      <Table striped bordered hover responsive>
        <thead>
          <tr>
            <th>Suite ID</th>
            <th>Name</th>
            <th>Created</th>
          </tr>
        </thead>
        <tbody>
          {suites.map((suite) => (
            <tr
              key={suite.id}
              onClick={() => handleSuiteClick(suite.id)}
              style={{ cursor: 'pointer' }}
              className="clickable-row"
            >
              <td>#{suite.id}</td>
              <td>{suite.name}</td>
              <td>{new Date(suite.time).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </Table>
      {suites.length === 0 && !loading && (
        <Alert variant="info" className="mt-3">No suites found for this project.</Alert>
      )}
    </div>
  );
};

export default SuitesTable;
