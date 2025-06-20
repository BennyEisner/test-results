import { useParams, useNavigate } from 'react-router-dom';
import { Container, Row, Col, Button, Alert } from 'react-bootstrap';
import BuildsTable from './BuildsTable';

const SuiteDetail = () => {
  const { suiteId, projectId } = useParams<{ suiteId: string; projectId: string }>();
  const navigate = useNavigate();

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

  return (
    <Container className="py-3 suite-detail">
      <Row className="align-items-center mb-3 suite-header">
        <Col xs="auto">
          <Button 
            variant="outline-secondary" 
            onClick={() => navigate(`/projects/${projectId}`)}
          >
            &laquo; Back to Project Suites
          </Button>
        </Col>
        <Col>
          <h1 className="h3 mb-0">Suite Details: #{suiteId} (Project: #{projectId})</h1>
        </Col>
      </Row>

      <BuildsTable projectId={projectId} suiteId={suiteId} />
    </Container>
  );
};

export default SuiteDetail;
