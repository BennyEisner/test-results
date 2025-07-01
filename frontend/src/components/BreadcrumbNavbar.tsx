import { useState, useEffect } from 'react';
import { useLocation, useParams, useNavigate } from 'react-router-dom';
import { fetchProjects } from '../services/api';
import type { Project, SearchResult } from '../types';
import { useDashboard } from '../context/DashboardContext';
import SearchBar from './SearchBar';
import './BreadcrumbNavbar.css';

interface BreadcrumbItem {
    label: string;
    path?: string;
}

const BreadcrumbNavbar = () => {
    const location = useLocation();
    const params = useParams();
    const navigate = useNavigate();
    const [breadcrumbs, setBreadcrumbs] = useState<BreadcrumbItem[]>([]);
    const [projects, setProjects] = useState<Project[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const dashboardContext = useDashboard();

    useEffect(() => {
        const loadProjects = async () => {
            try {
                const projectData = await fetchProjects();
                setProjects(projectData);
            } catch (error) {
                console.error('Failed to fetch projects:', error);
            }
        };

        loadProjects();
    }, []);

    useEffect(() => {
        const buildBreadcrumbs = async () => {
            setIsLoading(true);
            const items: BreadcrumbItem[] = [];

      try {
        if (location.pathname === '/') {
          items.push({ label: 'Dashboard' });
        } else {
          items.push({ label: 'Projects', path: '/projects' });
        }

        // Handle projectspecific routes
        if (params.projectId) {
                    const projects = await fetchProjects();
                    const project = projects.find(p => p.id.toString() === params.projectId);

                    if (project) {
                        items.push({
                            label: project.name,
                            path: `/projects/${params.projectId}`
                        });

                        // Handle suite-specific routes
                        if (params.suiteId) {
                            // For now, we'll use a placeholder for suite name
                            // In a real implementation, you'd fetch suite details
                            items.push({
                                label: `Suite ${params.suiteId}`,
                                path: `/projects/${params.projectId}/suites/${params.suiteId}`
                            });

                            // Handle builds route
                            if (location.pathname.includes('/builds')) {
                                if (params.buildId) {
                                    items.push({
                                        label: 'Builds',
                                        path: `/projects/${params.projectId}/suites/${params.suiteId}/builds`
                                    });
                                    items.push({
                                        label: `Build ${params.buildId}`
                                    });
                                } else {
                                    items.push({ label: 'Builds' });
                                }
                            }
                        }
                    }
                }

                setBreadcrumbs(items);
            } catch (error) {
                console.error('Error building breadcrumbs:', error);
                // Fallback breadcrumbs
                setBreadcrumbs([{ label: 'Dashboard', path: '/' }]);
            } finally {
                setIsLoading(false);
            }
        };

        buildBreadcrumbs();
    }, [location.pathname, params]);

    const handleSearchResultSelect = (result: SearchResult) => {
        // Navigate based on the search result type
        switch (result.type) {
            case 'project':
                navigate(`/projects/${result.id}`);
                break;
            case 'test_suite':
                navigate(result.url);
                break;
            case 'build':
                navigate(result.url);
                break;
            case 'test_case':
                navigate(result.url);
                break;
            default:
                console.log('Unknown result type:', result);
        }
    };

    const handleBreadcrumbClick = (path: string) => {
        navigate(path);
    };

    if (isLoading) {
        return (
            <div className="breadcrumb-navbar loading">
                <div className="breadcrumb-loading">Loading...</div>
            </div>
        );
    }

    const renderProjectSelector = () => {
        if (!dashboardContext) return null;

        return (
            <div className="project-selector-scroll-container">
                <div className="project-selector">
                    {projects.map((project) => (
                        <button
                            key={project.id}
                            className={`project-item ${dashboardContext.selectedProjectId === project.id ? 'active' : ''}`}
                            onClick={() => dashboardContext.onProjectSelect(project.id)}
                        >
                            {project.name}
                        </button>
                    ))}
                </div>
            </div>
        );
    };

    return (
        <div className="breadcrumb-navbar">
            <div className="breadcrumb-section">
                {dashboardContext && location.pathname.startsWith('/dashboard') ? (
                    renderProjectSelector()
                ) : (
                    breadcrumbs.map((item, index) => (
                        <span key={index} className="breadcrumb-item">
                            {item.path ? (
                                <button
                                    className="breadcrumb-link"
                                    onClick={() => handleBreadcrumbClick(item.path!)}
                                >
                                    {item.label}
                                </button>
                            ) : (
                                <span className="breadcrumb-current">{item.label}</span>
                            )}
                            {index < breadcrumbs.length - 1 && (
                                <span className="breadcrumb-separator"> {'>'} </span>
                            )}
                        </span>
                    ))
                )}
            </div>

            <div className="search-section">
                <SearchBar onResultSelect={handleSearchResultSelect} />
            </div>
        </div>
    );
};

export default BreadcrumbNavbar;
