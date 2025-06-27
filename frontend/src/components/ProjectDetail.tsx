{/*   /projects/{projectId}   */ }

import { useParams, useNavigate } from 'react-router-dom';
import { Button, Alert, Card } from 'react-bootstrap';
import SuitesTable from './SuitesTable.tsx';
import BuildsTable from './BuildsTable';

const ProjectDetail = () => {
    const { projectId } = useParams<{ projectId: string }>();
    const navigate = useNavigate();

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
                    onClick={() => navigate('/')}
                >
                    &laquo; Back to Dashboard
                </Button>
                <h1 className="page-title">Project: #{projectId}</h1>
            </div>
            <div className="grid-container">
                <Card className="overview-card">
                    <Card.Header as="h5">Suites</Card.Header>
                    <Card.Body>
                        <SuitesTable projectId={projectId} />
                    </Card.Body>
                </Card>
                <Card className="overview-card">
                    <Card.Header as="h5">Recent Builds</Card.Header>
                    <Card.Body>
                        <BuildsTable projectId={projectId} />
                    </Card.Body>
                </Card>
            </div>
        </div>
    );
};

export default ProjectDetail;
