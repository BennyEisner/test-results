import { useEffect } from 'react';
import { Card, Button, Container, Row, Col, Spinner, Alert } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { useAuth } from '../../context/AuthContext';
import './LoginPage.css';

const LoginPage = () => {
    const { login, isAuthenticated, isLoading, error, clearError } = useAuth();

    useEffect(() => {
        // Clear any previous errors when the component mounts
        return () => {
            clearError();
        };
    }, [clearError]);

    const handleGitHubLogin = () => {
        login('github');
    };

    const handleOktaLogin = () => {
        login('okta');
    };

    if (isAuthenticated) {
        return (
            <Container className="login-container">
                <Row className="justify-content-center">
                    <Col md={6} lg={4}>
                        <Card className="login-card text-center">
                            <Card.Body>
                                <Card.Title className="mb-4">
                                    <h2>You are already logged in</h2>
                                </Card.Title>
                                <Link to="/dashboard">
                                    <Button variant="primary" size="lg">
                                        Go to Dashboard
                                    </Button>
                                </Link>
                            </Card.Body>
                        </Card>
                    </Col>
                </Row>
            </Container>
        );
    }

    return (
        <Container className="login-container">
            <Row className="justify-content-center">
                <Col md={6} lg={4}>
                    <Card className="login-card">
                        <Card.Body className="text-center">
                            <Card.Title className="mb-4">
                                <h2>Welcome to Test Results</h2>
                            </Card.Title>

                            {error && <Alert variant="danger">{error}</Alert>}

                            <Card.Text className="mb-4">
                                Sign in to access your test results dashboard and manage your projects.
                            </Card.Text>

                            <div className="login-options">
                                <Button
                                    variant="outline-dark"
                                    size="lg"
                                    className="login-button mb-3"
                                    onClick={handleGitHubLogin}
                                    disabled={isLoading}
                                >
                                    {isLoading ? <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true" /> : <i className="fab fa-github me-2"></i>}
                                    Continue with GitHub
                                </Button>

                                <Button
                                    variant="primary"
                                    size="lg"
                                    className="login-button"
                                    onClick={handleOktaLogin}
                                    disabled={isLoading}
                                >
                                    {isLoading ? <Spinner as="span" animation="border" size="sm" role="status" aria-hidden="true" /> : <i className="fas fa-shield-alt me-2"></i>}
                                    Continue with Okta
                                </Button>
                            </div>

                            <div className="mt-4">
                                <small className="text-muted">
                                    Choose your preferred authentication method
                                </small>
                            </div>
                        </Card.Body>
                    </Card>
                </Col>
            </Row>
        </Container>
    );
};

export default LoginPage;
