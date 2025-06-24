import { useParams, useNavigate } from 'react-router-dom';
import { Container, Row, Col, Button, Alert, Card } from 'react-bootstrap';
import BuildsTable from './BuildsTable';
import TestCaseTrendChart from './TestCaseTrendChart';

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
    <Container fluid className="py-3 suite-detail" style={{ paddingLeft: '2rem', paddingRight: '2rem' }}>
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

      <Row>
        <Col md={6}>
          <Card className="mb-4">
            <Card.Body>
              <TestCaseTrendChart projectId={projectId} suiteId={suiteId} />
            </Card.Body>
          </Card>
        </Col>
        <Col md={6}>
          <BuildsTable projectId={projectId} suiteId={suiteId} />
        </Col>
      </Row>
    </Container>
  );
};

export default SuiteDetail;
