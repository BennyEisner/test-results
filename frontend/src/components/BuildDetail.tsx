import { useParams, useNavigate } from 'react-router-dom';
import { Button, Alert, Card, Row, Col } from 'react-bootstrap';
import ExecutionsTable from './ExecutionsTable';
import ExecutionsSummary from './ExecutionsSummary';
import BuildDoughnutChart from './BuildDoughnutChart';
import { useExecutionsSummary } from '../hooks/useExecutionsSummary';

const BuildDetail = () => {
    const { buildId, suiteId, projectId } = useParams<{ buildId: string; suiteId: string; projectId: string }>();
    const navigate = useNavigate();

    if (!suiteId) {
        return (
            <div className="page-container">
                <Alert variant="danger">Suite ID is required</Alert>
            </div>
        );
    }

    if (!projectId) {
        return (
            <div className="page-container">
                <Alert variant="danger">Project ID is required</Alert>
            </div>
        );
    }

    if (!buildId) {
        return (
            <div className="page-container">
                <Alert variant="danger">Build ID is required</Alert>
            </div>
        );
    }

    const { stats, executions, loading, error } = useExecutionsSummary(buildId);

    return (
        <div className="page-container">
            <div className="page-header">
                <Button
                    variant="outline-primary"
                    className="accent-button-outline"
                    onClick={() => navigate(`/projects/${projectId}/suites/${suiteId}`)}
                >
                    &laquo; Back to Builds
                </Button>
                <h1 className="page-title">Build: #{buildId}</h1>
            </div>

            <Row>
                <Col md={8}>
                    <Card className="overview-card mb-4">
                        <Card.Header as="h5">Executions Summary</Card.Header>
                        <Card.Body>
                            <ExecutionsSummary stats={stats} loading={loading} />
                        </Card.Body>
                    </Card>
                </Col>
                <Col md={4}>
                    <Card className="overview-card mb-4">
                        <Card.Header as="h5">Executions Chart</Card.Header>
                        <Card.Body>
                            <BuildDoughnutChart buildId={buildId} />
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            {error && (
                <Alert variant="danger" className="mt-3">
                    Error loading executions: {error}
                </Alert>
            )}

            <Card className="overview-card">
                <Card.Header as="h5">Executions</Card.Header>
                <Card.Body>
                    {!error && <ExecutionsTable executions={executions} loading={loading} />}
                </Card.Body>
            </Card>
        </div>
    );
};

export default BuildDetail;
