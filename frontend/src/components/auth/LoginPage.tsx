import React from 'react';
import { Card, Button, Container, Row, Col } from 'react-bootstrap';
import { useAuth } from '../../context/AuthContext';
import './LoginPage.css';

const LoginPage: React.FC = () => {
  const { login } = useAuth();

  const handleGitHubLogin = () => {
    login('github');
  };

  const handleOktaLogin = () => {
    login('okta');
  };

  return (
    <Container className="login-container">
      <Row className="justify-content-center">
        <Col md={6} lg={4}>
          <Card className="login-card">
            <Card.Body className="text-center">
              <Card.Title className="mb-4">
                <h2>Welcome to Test Results</h2>
              </Card.Title>
              
              <Card.Text className="mb-4">
                Sign in to access your test results dashboard and manage your projects.
              </Card.Text>

              <div className="login-options">
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