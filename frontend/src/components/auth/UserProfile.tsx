import React, { useState, useEffect } from 'react';
import { Card, Button, Table, Modal, Form, Alert, Badge } from 'react-bootstrap';
import { useAuth } from '../../context/AuthContext';
import { authApi } from '../../services/authApi';
import { User, APIKey } from '../../types/auth';
import './UserProfile.css';

const UserProfile: React.FC = () => {
  const { user, logout } = useAuth();
  const [apiKeys, setApiKeys] = useState<APIKey[]>([]);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [newKeyName, setNewKeyName] = useState('');
  const [newApiKey, setNewApiKey] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadAPIKeys();
  }, []);

  const loadAPIKeys = async () => {
    try {
      const keys = await authApi.listAPIKeys();
      setApiKeys(keys);
    } catch (err) {
      setError('Failed to load API keys');
    }
  };

  const handleCreateAPIKey = async () => {
    if (!newKeyName.trim()) {
      setError('Please enter a name for the API key');
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await authApi.createAPIKey({ name: newKeyName });
      setNewApiKey(response.plain_text_key);
      setApiKeys([response.api_key, ...apiKeys]);
      setNewKeyName('');
      setShowCreateModal(false);
    } catch (err) {
      setError('Failed to create API key');
    } finally {
      setLoading(false);
    }
  };

  const handleDeleteAPIKey = async (keyId: number) => {
    if (!window.confirm('Are you sure you want to delete this API key?')) {
      return;
    }

    try {
      await authApi.deleteAPIKey(keyId);
      setApiKeys(apiKeys.filter(key => key.id !== keyId));
    } catch (err) {
      setError('Failed to delete API key');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const isExpired = (expiresAt: string) => {
    return new Date(expiresAt) < new Date();
  };

  return (
    <div className="user-profile">
      <div className="container">
        <div className="row">
          <div className="col-md-4">
            <Card className="profile-card">
              <Card.Body className="text-center">
                {user?.avatar_url && (
                  <img 
                    src={user.avatar_url} 
                    alt="Avatar" 
                    className="profile-avatar mb-3"
                  />
                )}
                <Card.Title>{user?.name || 'User'}</Card.Title>
                <Card.Text className="text-muted">{user?.email}</Card.Text>
                <Badge bg="secondary" className="mb-3">
                  {user?.provider}
                </Badge>
                <div className="d-grid">
                  <Button variant="outline-danger" onClick={logout}>
                    Logout
                  </Button>
                </div>
              </Card.Body>
            </Card>
          </div>

          <div className="col-md-8">
            <Card>
              <Card.Header className="d-flex justify-content-between align-items-center">
                <h5 className="mb-0">API Keys</h5>
                <Button 
                  variant="primary" 
                  size="sm"
                  onClick={() => setShowCreateModal(true)}
                >
                  Create API Key
                </Button>
              </Card.Header>
              <Card.Body>
                {error && (
                  <Alert variant="danger" dismissible onClose={() => setError(null)}>
                    {error}
                  </Alert>
                )}

                {apiKeys.length === 0 ? (
                  <p className="text-muted">No API keys found. Create one to use with CLI or Jenkins.</p>
                ) : (
                  <Table responsive>
                    <thead>
                      <tr>
                        <th>Name</th>
                        <th>Created</th>
                        <th>Last Used</th>
                        <th>Expires</th>
                        <th>Status</th>
                        <th>Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {apiKeys.map(key => (
                        <tr key={key.id}>
                          <td>{key.name}</td>
                          <td>{formatDate(key.created_at)}</td>
                          <td>
                            {key.last_used_at 
                              ? formatDate(key.last_used_at)
                              : 'Never'
                            }
                          </td>
                          <td>{formatDate(key.expires_at)}</td>
                          <td>
                            <Badge bg={isExpired(key.expires_at) ? 'danger' : 'success'}>
                              {isExpired(key.expires_at) ? 'Expired' : 'Active'}
                            </Badge>
                          </td>
                          <td>
                            <Button 
                              variant="outline-danger" 
                              size="sm"
                              onClick={() => handleDeleteAPIKey(key.id)}
                              disabled={isExpired(key.expires_at)}
                            >
                              Delete
                            </Button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </Table>
                )}
              </Card.Body>
            </Card>
          </div>
        </div>
      </div>

      {/* Create API Key Modal */}
      <Modal show={showCreateModal} onHide={() => setShowCreateModal(false)}>
        <Modal.Header closeButton>
          <Modal.Title>Create API Key</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Form>
            <Form.Group>
              <Form.Label>API Key Name</Form.Label>
              <Form.Control
                type="text"
                placeholder="e.g., Jenkins CLI, Development"
                value={newKeyName}
                onChange={(e) => setNewKeyName(e.target.value)}
              />
              <Form.Text className="text-muted">
                Give your API key a descriptive name to help you identify its purpose.
              </Form.Text>
            </Form.Group>
          </Form>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="secondary" onClick={() => setShowCreateModal(false)}>
            Cancel
          </Button>
          <Button 
            variant="primary" 
            onClick={handleCreateAPIKey}
            disabled={loading || !newKeyName.trim()}
          >
            {loading ? 'Creating...' : 'Create API Key'}
          </Button>
        </Modal.Footer>
      </Modal>

      {/* New API Key Display Modal */}
      <Modal show={!!newApiKey} onHide={() => setNewApiKey(null)}>
        <Modal.Header closeButton>
          <Modal.Title>API Key Created</Modal.Title>
        </Modal.Header>
        <Modal.Body>
          <Alert variant="warning">
            <strong>Important:</strong> This is the only time you'll see this API key. 
            Make sure to copy it now and store it securely.
          </Alert>
          <Form.Group>
            <Form.Label>Your API Key</Form.Label>
            <Form.Control
              as="textarea"
              rows={3}
              value={newApiKey || ''}
              readOnly
            />
          </Form.Group>
        </Modal.Body>
        <Modal.Footer>
          <Button variant="primary" onClick={() => setNewApiKey(null)}>
            I've Copied It
          </Button>
        </Modal.Footer>
      </Modal>
    </div>
  );
};

export default UserProfile; 