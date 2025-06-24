{/*   /projects/{projectId}   */ }

import { useParams, useNavigate } from 'react-router-dom';
import { Container, Row, Col, Button, Alert } from 'react-bootstrap';
import SuitesTable from './SuitesTable.tsx';

const ProjectDetail = () => {
    const { projectId } = useParams<{ projectId: string }>();
    const navigate = useNavigate();

    if (!projectId) {
        return (
            <Container className="py-3">
                <Alert variant="danger">Project ID is required</Alert>
            </Container>
        );
    }

    return (
        <Container className="py-3 project-detail">
            <Row className="align-items-center mb-3 project-header">
                <Col xs="auto">
                    <Button
                        variant="outline-secondary"
                        onClick={() => navigate('/')}
                    >
                        &laquo; Back to Dashboard
                    </Button>
                </Col>
                <Col>
                    <h1 className="h3 mb-0">Project Details: #{projectId}</h1>
                </Col>
            </Row>

            <SuitesTable projectId={projectId} />
        </Container>
    );
};

export default ProjectDetail;
