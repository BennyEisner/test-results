import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Table, Spinner, Alert, Badge } from 'react-bootstrap';
import type { Build } from '../types';
import { fetchBuilds } from '../services/api';

interface BuildsTableProps {
    projectId?: string | number;
    suiteId?: string | number;
    fetchFunction?: () => Promise<Build[]>;
    title?: string;
}

const BuildsTable = ({ projectId, suiteId, fetchFunction, title }: BuildsTableProps) => {
    const [builds, setBuilds] = useState<Build[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const navigate = useNavigate(); // Initialize navigate

    useEffect(() => {
        const loadBuilds = async () => {
            try {
                setLoading(true);
                let data;
                // Handle component based on props
                if (fetchFunction) {
                    data = await fetchFunction();
                } else if (projectId && suiteId) {
                    // Pass projectId and suiteId to fetchBuilds
                    data = await fetchBuilds(projectId, suiteId);
                }
                else {
                    throw new Error("Either fetchFunction or projectId and suiteId must be specified");
                }
                setBuilds(data);
                setError(null);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load builds');
            } finally {
                setLoading(false);
            }
        };

        loadBuilds();
    }, [projectId, suiteId, fetchFunction]);

    const handleBuildClick = (build: Build) => {
        // Navigate to the executions page for the clicked build
        console.log('Build object:', build);
        console.log('projectId prop:', projectId);
        console.log('build.project_id:', build.project_id);
        const targetProjectId = projectId ?? build.project_id;
        const targetSuiteId = suiteId ?? build.test_suite_id;
        console.log('targetProjectId:', targetProjectId);
        console.log('targetSuiteId:', targetSuiteId);
        navigate(`/projects/${targetProjectId}/suites/${targetSuiteId}/builds/${build.id}`);
    };

    if (loading) {
        return (
            <div className="d-flex justify-content-center align-items-center" style={{ height: '80vh' }}>
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading builds...</span>
                </Spinner>
            </div>
        );
    }

    if (error) {
        return <Alert variant="danger">Error: {error}</Alert>;
    }

    return (
        <div className="py-3">
            {title && <h2 className="mb-3">{title}</h2>}
            <Table striped bordered hover responsive>
                <thead>
                    <tr>
                        <th>Build ID</th>
                        <th>Build Number</th>
                        <th>Test Cases</th>
                        <th>CI Provider</th>
                        <th>Created</th>
                    </tr>
                </thead>
                <tbody>
                    {builds.map((build) => (
                        <tr key={build.id} onClick={() => handleBuildClick(build)} style={{ cursor: 'pointer' }}>
                            <td>#{build.id}</td>
                            <td className="font-monospace">{build.build_number}</td>
                            <td>{build.test_case_count}</td>
                            <td>
                                <Badge bg={getCIProviderBadgeColor(build.ci_provider)}>{build.ci_provider || 'N/A'}</Badge>
                            </td>
                            <td className="text-muted fst-italic">{new Date(build.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </Table>
            {builds.length === 0 && !loading && (
                <Alert variant="info" className="mt-3">No builds found for this project.</Alert>
            )}
        </div>
    );
};

// Helper function to determine badge color based on CI provider
const getCIProviderBadgeColor = (provider: string | null | undefined) => {
    if (!provider) return 'secondary';
    const lowerProvider = provider.toLowerCase();
    if (lowerProvider.includes('github')) return 'dark';
    if (lowerProvider.includes('jenkins')) return 'danger';
    if (lowerProvider.includes('travis')) return 'info';
    return 'secondary';
};

export default BuildsTable;
