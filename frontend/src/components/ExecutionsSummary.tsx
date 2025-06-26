// ExecutionsSummary.tsx
import { Card, Row, Col, Spinner } from 'react-bootstrap';
import { useExecutionsSummary } from '../hooks/useExecutionsSummary';

interface ExecutionsSummaryProps {
    buildId?: string | number;
    title?: string;
}

const ExecutionsSummary = ({ buildId, title }: ExecutionsSummaryProps) => {
    if (!buildId) {
        return <div className="text-muted">Please provide a build ID.</div>;
    }

    const { stats, loading } = useExecutionsSummary(buildId);

    if (loading) {
        return (
            <div className="d-flex justify-content-center align-items-center my-3">
                <Spinner animation="border" role="status" variant="primary" />
                <span className="ms-2">Analyzing test results...</span>
            </div>
        );
    }

    if (stats.total === 0) {
        return (
            <div className="my-3 text-center text-muted">
                No execution data to summarize.
            </div>
        );
    }

    // Render the summary metrics
    return (
        <div>
            {title && <h3 className="component-title">{title}</h3>}
            <Row xs={1} sm={2} md={3} lg={6} className="g-3">
                <Col>
                    <Card text="white" bg="primary" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.total}</Card.Title>
                            <Card.Text>Total Tests</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
                <Col>
                    <Card text="white" bg="success" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.passed}</Card.Title>
                            <Card.Text>Passed</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
                <Col>
                    <Card text="white" bg="danger" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.failed}</Card.Title>
                            <Card.Text>Failed</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
                <Col>
                    <Card text="white" bg="secondary" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.skipped}</Card.Title>
                            <Card.Text>Skipped</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
                <Col>
                    <Card text="white" bg="primary" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.passRate}%</Card.Title>
                            <Card.Text>Pass Rate</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
                <Col>
                    <Card text="white" bg="primary" className="h-100">
                        <Card.Body className="text-center">
                            <Card.Title as="h4">{stats.avgTime.toFixed(2)}s</Card.Title>
                            <Card.Text>Avg. Time</Card.Text>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default ExecutionsSummary;
