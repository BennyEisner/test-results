import { useState, useEffect } from 'react';
import { Nav, Navbar, Container, NavDropdown, Button } from 'react-bootstrap';
import { Link, useNavigate } from 'react-router-dom';
import { fetchProjects } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import type { Project } from '../../types';
import './AppNavbar.css';

interface AppNavbarProps {
    onProjectClick?: (projectId: number) => void;
}

const AppNavbar = ({ onProjectClick }: AppNavbarProps) => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [error, setError] = useState<string | null>(null);
    const { user, isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();

    useEffect(() => {
        const loadProjects = async () => {
            try {
                const data = await fetchProjects();
                setProjects(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load projects in navbar');
                console.error("Navbar project fetch error: ", err);
            }
        };
        loadProjects();
    }, []);

    const handleLogout = async () => {
        await logout();
        navigate('/');
    };

    if (error) {
        return (
            <div className="navbar-error">
                <span>Error loading projects: {error}</span>
            </div>
        );
    }

    return (
        <Navbar bg="light" expand="lg" className="app-navbar">
            <Container>
                <Navbar.Brand as={Link} to="/">
                    Test Results
                </Navbar.Brand>
                
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        {isAuthenticated && (
                            <>
                                <Nav.Link as={Link} to="/dashboard">
                                    Dashboard
                                </Nav.Link>
                                <Nav.Link as={Link} to="/projects">
                                    Projects
                                </Nav.Link>
                            </>
                        )}
                    </Nav>
                    
                    <Nav className="projects-nav-container">
                        {isAuthenticated && projects.map((project) => {
                            if (onProjectClick) {
                                return (
                                    <Nav.Link
                                        key={project.id}
                                        onClick={() => onProjectClick(project.id)}
                                        className="project-nav-link"
                                        style={{ cursor: 'pointer' }}
                                    >
                                        {project.name}
                                    </Nav.Link>
                                );
                            }
                            return (
                                <Nav.Link
                                    key={project.id}
                                    as={Link}
                                    to={`/projects/${project.id}`}
                                    className="project-nav-link"
                                >
                                    {project.name}
                                </Nav.Link>
                            );
                        })}
                    </Nav>

                    <Nav className="ms-auto">
                        {isAuthenticated ? (
                            <NavDropdown 
                                title={
                                    <span>
                                        {user?.avatar_url && (
                                            <img 
                                                src={user.avatar_url} 
                                                alt="Avatar" 
                                                className="navbar-avatar me-2"
                                            />
                                        )}
                                        {user?.name || 'User'}
                                    </span>
                                } 
                                id="user-dropdown"
                            >
                                <NavDropdown.Item as={Link} to="/profile">
                                    Profile & API Keys
                                </NavDropdown.Item>
                                <NavDropdown.Divider />
                                <NavDropdown.Item onClick={handleLogout}>
                                    Logout
                                </NavDropdown.Item>
                            </NavDropdown>
                        ) : (
                            <Nav.Link as={Link} to="/">
                                Login
                            </Nav.Link>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
};

export default AppNavbar;
