import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Button, Alert, Card } from 'react-bootstrap';
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

            <Card className="overview-card mb-4">
                <Card.Header as="h5">Executions Summary</Card.Header>
                <Card.Body>
                    <ExecutionsSummary executions={executions} loading={loading} />
                </Card.Body>
            </Card>

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
