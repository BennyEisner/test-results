import { Row, Col, Card, Alert } from 'react-bootstrap';
import AppNavbar from './AppNavbar';
import BuildsTable from './BuildsTable';
import { fetchRecentBuilds } from '../services/api';
import './HomePage.css';

const HomePage = () => {
    return (
        <div >
            {/* Projects Navigation Card */}
            <Row className="mb-4">
                <Col>
 
                            <AppNavbar />
  
                </Col>
            </Row>

            <Row>
                {/* Recent Builds Card */}
                <Col md={12} className="mb-4">
                    <Card className="overview-card">
                        <Card.Header as="h5">Recent Builds</Card.Header>
                        <Card.Body>
                            <div className="builds-table-container">
                                <BuildsTable fetchFunction={fetchRecentBuilds} title="" />
                            </div>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>

            <Row>
                {/* Global Statistics Card */}
                <Col md={6} className="mb-4">
                    <Card className="overview-card">
                        <Card.Header as="h5">Global Statistics</Card.Header>
                        <Card.Body>
                            <Alert variant="info" className="info-alert">
                                Placeholder for global statistics (e.g., total projects, overall pass rate).
                            </Alert>
                        </Card.Body>
                    </Card>
                </Col>

                {/* Recent Failures Card */}
                <Col md={6} className="mb-4">
                    <Card className="overview-card">
                        <Card.Header as="h5">Recent Failures</Card.Header>
                        <Card.Body>
                            <Alert variant="info" className="info-alert">
                                Placeholder for a panel showing recent failures.
                            </Alert>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default HomePage;
