{/*   /projects/{projectId}   */ }

import { useParams, useNavigate } from 'react-router-dom';
import { Button, Alert } from 'react-bootstrap';
import SuitesTable from './SuitesTable.tsx';

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
            <SuitesTable projectId={projectId} />
        </div>
    );
};

export default ProjectDetail;
