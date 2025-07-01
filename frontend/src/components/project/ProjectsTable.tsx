{/*   /projects   */ }

import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Table, Spinner, Alert, Card, Container } from 'react-bootstrap';
import type { Project } from '../../types';
import { fetchProjects } from '../../services/api';

const ProjectsTable = () => {
    const navigate = useNavigate();
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadProjects = async () => {
            try {
                setLoading(true);
                const data = await fetchProjects();
                setProjects(data);
                setError(null);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load projects');
            } finally {
                setLoading(false);
            }
        };
        loadProjects();
    }, []);

    const handleProjectClick = (projectId: string | number) => {
        navigate(`/projects/${projectId}`);
    };

    if (loading) {
        return (
            <div className="page-container">
                <div className="d-flex justify-content-center align-items-center" style={{ height: '80vh' }}>
                    <Spinner animation="border" role="status">
                        <span className="visually-hidden">Loading projects...</span>
                    </Spinner>
                </div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="page-container">
                <Container fluid>
                    <Alert variant="danger">Error: {error}</Alert>
                </Container>
            </div>
        );
    }

    return (
        <div className="page-container">
            <Container fluid>
                <h1 className="page-title mb-4">Projects</h1>
                <Card className="overview-card">
                    <Card.Header as="h5">All Projects</Card.Header>
                    <Card.Body>
                        <Table bordered hover responsive>
                            <thead>
                                <tr>
                                    <th>ID</th>
                                    <th>Project Name</th>
                                </tr>
                            </thead>
                            <tbody>
                                {projects.map((project) => (
                                    <tr
                                        key={project.id}
                                        onClick={() => handleProjectClick(project.id)}
                                        style={{ cursor: 'pointer' }}
                                        className="clickable-row"
                                    >
                                        <td>{project.id}</td>
                                        <td>{project.name}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </Table>
                        {projects.length === 0 && !loading && (
                            <Alert variant="info" className="info-alert mt-3">No projects found.</Alert>
                        )}
                    </Card.Body>
                </Card>
            </Container>
        </div>
    );
};

export default ProjectsTable;
