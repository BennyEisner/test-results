import { useState, useEffect } from 'react';
import { Nav } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { fetchProjects } from '../../services/api';
import type { Project } from '../../types';
import './AppNavbar.css';

interface AppNavbarProps {
    onProjectClick?: (projectId: number) => void;
}

const AppNavbar = ({ onProjectClick }: AppNavbarProps) => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [error, setError] = useState<string | null>(null);

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

    if (error) {
        return (
            <div className="navbar-error">
                <span>Error loading projects: {error}</span>
            </div>
        );
    }

    return (
        <div className="app-navbar-container">
            <Nav className="projects-nav-container">
                {projects.map((project) => {
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
        </div>
    );
};

export default AppNavbar;
