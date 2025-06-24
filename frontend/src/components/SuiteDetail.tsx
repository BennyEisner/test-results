import { useParams, useNavigate } from 'react-router-dom';
import { Button, Alert, Card } from 'react-bootstrap';
import BuildsTable from './BuildsTable';
import TestCaseTrendChart from './TestCaseTrendChart';

const SuiteDetail = () => {
    const { suiteId, projectId } = useParams<{ suiteId: string; projectId: string }>();
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

    return (
        <div className="page-container">
            <div className="page-header">
                <Button
                    variant="outline-primary"
                    className="accent-button-outline"
                    onClick={() => navigate(`/projects/${projectId}`)}
                >
                    &laquo; Back to Project Suites
                </Button>
                <h1 className="page-title">Suite: #{suiteId}</h1>
            </div>
            <div className="grid-container">
                <Card className="overview-card">
                    <Card.Header as="h5">Test Case Trend</Card.Header>
                    <Card.Body>
                        <TestCaseTrendChart projectId={projectId} suiteId={suiteId} />
                    </Card.Body>
                </Card>
                <Card className="overview-card">
                    <Card.Header as="h5">Builds</Card.Header>
                    <Card.Body>
                        <BuildsTable projectId={projectId} suiteId={suiteId} />
                    </Card.Body>
                </Card>
            </div>
        </div>
    );
};

export default SuiteDetail;
