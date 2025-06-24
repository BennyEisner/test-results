{/*   /   */ }

import { Container, Row, Col, Card, Alert, Button } from 'react-bootstrap';
import { Link } from 'react-router-dom'
import AppNavbar from './AppNavbar';
import BuildsTable from './BuildsTable';
import { fetchRecentBuilds } from '../services/api';

const HomePage = () => {
    return (
        <Container fluid className="py-3">
            <Row className="mb-3">
                <Col>
                    <h1>Test Results Dashboard Overview</h1>
                    <p className="lead">Test Results Dashboard</p>
                </Col>
            </Row>

            <Row className="mb-4">
                <Col>
                    <Card className="shadow-sm">
                        <Card.Header as="h5" className="bg-light">Projects</Card.Header>
                        <AppNavbar />
                    </Card>
                </Col>
            </Row>

            {/* Section 1: Placeholder for Global Stats */}
            <Row className="mb-4">
                <Col>
                    <Card className="shadow-sm">
                        <Card.Header as="h5" className="bg-light">Global Statistics</Card.Header>
                        <Card.Body>
                            <Alert variant="info">
                                Placeholder for global statistics (e.g., total projects, overall pass rate).
                            </Alert>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            {/* Section 2: Placeholder for Recent Activity / Failures */}
            <Row className="mb-4">
                <Col md={6} className="mb-3 mb-md-0">
                    <Card className="shadow-sm" style={{ height: '400px' }}>
                        <Card.Header as="h5" className="bg-light">Recent Builds</Card.Header>
                        <Card.Body className="p-0" style={{ overflow: 'hidden' }}>
                            <div style={{ height: '100%', overflowY: 'auto' }}>
                                <BuildsTable
                                    fetchFunction={fetchRecentBuilds}
                                    title=""
                                />
                            </div>
                        </Card.Body>
                    </Card>
                </Col>
                <Col md={6}>
                    <Card className="shadow-sm">
                        <Card.Header as="h5" className="bg-light">Recent Failures</Card.Header>
                        <Card.Body>
                            <Alert variant="info">
                                Placeholder for a panel showing recent failures.
                            </Alert>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            {/* Section 3: Placeholder for Project Overview or Search */}
            <Row>
                <Col>
                    <Card className="shadow-sm">
                        <Card.Header as="h5" className="bg-light">Projects Overview</Card.Header>
                        <Card.Body>
                            <Alert variant="info">
                            </Alert>
                            <Link to="/projects">
                                <Button variant="outline-primary" size="sm" className="mt-2">
                                    View All Projects
                                </Button>
                            </Link>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Container>
    );
};

export default HomePage;
