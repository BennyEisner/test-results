import { Row, Col, Card, Alert, Button, Container } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import AppNavbar from '../common/AppNavbar';
import BuildsTable from '../build/BuildsTable';
import { useAuth } from '../../context/AuthContext';
import { authApi } from '../../services/authApi';
import './HomePage.css';
import { fetchBuilds } from '../../services/api';

const HomePage = () => {
    const { isAuthenticated, user } = useAuth();

    const handleGitHubLogin = () => {
        authApi.beginOAuth2Auth('github');
    };

    const handleOktaLogin = () => {
        authApi.beginOAuth2Auth('okta');
    };

    if (!isAuthenticated) {
        return (
            <div>
                <AppNavbar />
                <Container className="welcome-container">
                    <Row className="justify-content-center">
                        <Col md={8} lg={6}>
                            <Card className="welcome-card text-center">
                                <Card.Body className="p-5">
                                    <h1 className="mb-4">Welcome to Test Results</h1>
                                    <p className="lead mb-4">
                                        Track, analyze, and visualize your test results across all your projects. 
                                        Get insights into test performance and identify areas for improvement.
                                    </p>
                                    
                                    <div className="login-options mb-4">
                                        <Button 
                                            variant="outline-dark" 
                                            size="lg" 
                                            className="login-button mb-3"
                                            onClick={handleGitHubLogin}
                                        >
                                            <i className="fab fa-github me-2"></i>
                                            Continue with GitHub
                                        </Button>

                                        <Button 
                                            variant="primary" 
                                            size="lg" 
                                            className="login-button"
                                            onClick={handleOktaLogin}
                                        >
                                            <i className="fas fa-shield-alt me-2"></i>
                                            Continue with Okta
                                        </Button>
                                    </div>

                                    <div className="features">
                                        <Row>
                                            <Col md={4}>
                                                <div className="feature-item">
                                                    <i className="fas fa-chart-line fa-2x text-primary mb-2"></i>
                                                    <h5>Analytics</h5>
                                                    <p>Comprehensive test analytics and trends</p>
                                                </div>
                                            </Col>
                                            <Col md={4}>
                                                <div className="feature-item">
                                                    <i className="fas fa-project-diagram fa-2x text-primary mb-2"></i>
                                                    <h5>Projects</h5>
                                                    <p>Organize tests by projects and suites</p>
                                                </div>
                                            </Col>
                                            <Col md={4}>
                                                <div className="feature-item">
                                                    <i className="fas fa-key fa-2x text-primary mb-2"></i>
                                                    <h5>API Access</h5>
                                                    <p>CLI and Jenkins integration support</p>
                                                </div>
                                            </Col>
                                        </Row>
                                    </div>
                                </Card.Body>
                            </Card>
                        </Col>
                    </Row>
                </Container>
            </div>
        );
    }

    return (
        <div>
            <AppNavbar />
            
            {/* Welcome message for authenticated users */}
            <Row className="mb-4">
                <Col>
                    <Alert variant="success" className="welcome-alert">
                        <h5>Welcome back, {user?.name || 'User'}!</h5>
                        <p className="mb-0">
                            Ready to dive into your test results? Check out your{' '}
                            <Link to="/dashboard">dashboard</Link> or explore your{' '}
                            <Link to="/projects">projects</Link>.
                        </p>
                    </Alert>
                </Col>
            </Row>

            <Row>
                {/* Recent Builds Card */}
                <Col md={12} className="mb-4">
                    <Card className="overview-card">
                        <Card.Header as="h5">Recent Builds</Card.Header>
                        <Card.Body>
                            <div className="builds-table-container">
                                <BuildsTable fetchFunction={() => fetchBuilds(1)} title="" />
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
