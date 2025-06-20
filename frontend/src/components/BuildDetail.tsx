import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Container, Row, Col, Button, Alert } from 'react-bootstrap';
import ExecutionsTable from './ExecutionsTable';
import ExecutionsSummary from './ExecutionsSummary';
import { fetchExecutions } from '../services/api';
import type { TestCaseExecution } from '../types';

const BuildDetail = () => {
  const { buildId, suiteId, projectId } = useParams<{ buildId: string; suiteId: string; projectId: string }>();
  const navigate = useNavigate();

  const [executions, setExecutions] = useState<TestCaseExecution[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  if (!suiteId) {
    return (
      <Container className="py-3">
        <Alert variant="danger">Suite ID is required</Alert>
      </Container>
    );
  }

  if (!projectId) {
    return (
      <Container className="py-3">
        <Alert variant="danger">Project ID is required</Alert>
      </Container>
    );
  }

  if (!buildId) {
    return (
      <Container className="py-3">
        <Alert variant="danger">Build ID is required</Alert>
      </Container>
    );
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
    <Container className="py-3 build-detail">
      <Row className="align-items-center mb-3 build-header">
        <Col xs="auto">
          <Button 
            variant="outline-secondary" 
            onClick={() => navigate(`/projects/${projectId}/suites/${suiteId}`)}
          >
            &laquo; Back to Builds
          </Button>
        </Col>
        <Col>
          <h1 className="h3 mb-0">Build Details: #{buildId} (Suite: #{suiteId}, Project: #{projectId})</h1>
        </Col>
      </Row>

      <ExecutionsSummary executions={executions} loading={loading} />

      {error && (
        <Alert variant="danger" className="mt-3">
          Error loading executions: {error}
        </Alert>
      )}
      
      {/* ExecutionsTable will be refactored separately. 
          It already handles its own loading/error state internally if not passed down.
          For now, we only render it if there's no top-level error from fetchExecutions. 
      */}
      {!error && <ExecutionsTable executions={executions} loading={loading} />}
    </Container>
  );
};

export default BuildDetail;
