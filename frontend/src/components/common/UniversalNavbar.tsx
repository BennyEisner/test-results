import { useState, useEffect } from 'react';
import { Nav, Navbar, Container, NavDropdown } from 'react-bootstrap';
import { Link, useNavigate, useLocation } from 'react-router-dom';
import { fetchProjects } from '../../services/api';
import { useAuth } from '../../context/AuthContext';
import type { Project, SearchResult } from '../../types';
import './UniversalNavbar.scss';
import SearchBar from './SearchBar';

interface UniversalNavbarProps {
    onProjectSelect?: (projectId: number) => void;
}

const UniversalNavbar = ({ onProjectSelect }: UniversalNavbarProps) => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [error, setError] = useState<string | null>(null);
    const { user, isAuthenticated, logout } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

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
        if (isAuthenticated) {
            loadProjects();
        }
    }, [isAuthenticated]);

    const handleLogout = async () => {
        await logout();
        navigate('/');
    };

    const handleResultSelect = (result: SearchResult) => {
        const { type, id, project_id } = result;
        switch (type) {
            case 'project':
                navigate(`/projects/${id}`);
                break;
            case 'test_suite':
                navigate(`/projects/${project_id}/suites/${id}`);
                break;
            case 'build':
                navigate(`/builds/${id}`);
                break;
            case 'test_case':
                navigate(`/test-cases/${id}`);
                break;
            default:
                console.warn(`Unknown search result type: ${type}`);
        }
    };

    const renderTitle = () => {
        if (location.pathname === '/') {
            return "Test Results";
        }
        if (location.pathname.startsWith('/dashboard')) {
            return "Dashboard";
        }
        if (location.pathname.startsWith('/projects')) {
            return "Projects";
        }
        return "Test Results";
    };

    const renderNavLinks = () => {
        if (location.pathname === '/') {
            return (
                <>
                    <Nav.Link as={Link} to="/dashboard">
                        Dashboard
                    </Nav.Link>
                    <Nav.Link as={Link} to="/projects">
                        Projects
                    </Nav.Link>
                </>
            );
        }
        if (location.pathname.startsWith('/dashboard')) {
            return (
                <>
                    <Nav.Link as={Link} to="/">
                        Home
                    </Nav.Link>
                    <NavDropdown title="Projects" id="projects-dropdown">
                        {projects.map((project) => (
                            <NavDropdown.Item
                                key={project.id}
                                onClick={() => onProjectSelect && onProjectSelect(project.id)}
                            >
                                {project.name}
                            </NavDropdown.Item>
                        ))}
                    </NavDropdown>
                </>
            );
        }
        return (
            <Nav.Link as={Link} to="/">
                Home
            </Nav.Link>
        );
    };

    if (error) {
        return (
            <div className="navbar-error">
                <span>Error loading projects: {error}</span>
            </div>
        );
    }

    return (
        <Navbar expand="lg" className="app-navbar">
            <Container>
                <Navbar.Brand as={Link} to="/">
                    {renderTitle()}
                </Navbar.Brand>
                
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="me-auto">
                        {isAuthenticated && renderNavLinks()}
                    </Nav>
                    
                    <Nav className="ms-auto d-flex align-items-center">
                        <SearchBar onResultSelect={handleResultSelect} />
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
                                align="end"
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

export default UniversalNavbar;
