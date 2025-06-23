import { useState, useEffect } from 'react';
import { Nav } from 'react-bootstrap';
import { Link } from 'react-router-dom';
import { fetchProjects } from '../services/api';
import type { Project } from '../types';

const AppNavbar = () => {
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
            <div className="bg-light p-3 border-bottom">
                <span className="text-danger">Error loading projects: {error}</span>
            </div>
        );
    }

    return (
        <div className="p-2">
            <div 
                className="d-flex overflow-auto" 
                style={{ 
                    whiteSpace: 'nowrap',
                    scrollbarWidth: 'none', /* Firefox */
                }}
            >
                <style>
                    {`
                        .d-flex.overflow-auto::-webkit-scrollbar {
                            display: none; /* Safari and Chrome */
                        }
                    `}
                </style>
                <Nav className="flex-nowrap">
                    {projects.map((project) => (
                        <Nav.Link 
                            key={project.id} 
                            as={Link} 
                            to={`/projects/${project.id}`}
                            className="text-primary px-2 py-1 me-2 bg-white rounded border"
                            style={{ 
                                textDecoration: 'none',
                                minWidth: 'fit-content',
                                whiteSpace: 'nowrap',
                                boxShadow: '0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24)'
                            }}
                        >
                            {project.name}
                        </Nav.Link>
                    ))}
                </Nav>
            </div>
        </div>
    );
};

export default AppNavbar;
