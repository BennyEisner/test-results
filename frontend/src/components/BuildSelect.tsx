import { useState, useEffect } from 'react';
import { Spinner, Alert } from 'react-bootstrap';
import type { Build } from '../types';
import { fetchBuilds } from '../services/api';
import './BuildSelect.css';

interface BuildSelectProps {
    projectId: string | number;
    suiteId: string | number;
    onBuildSelect: (buildId: string | number) => void;
}

const BuildSelect = ({ projectId, suiteId, onBuildSelect }: BuildSelectProps) => {
    const [builds, setBuilds] = useState<Build[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedBuildId, setSelectedBuildId] = useState<string | number | null>(null);

    useEffect(() => {
        const loadBuilds = async () => {
            try {
                setLoading(true);
                const data = await fetchBuilds(projectId, suiteId);
                setBuilds(data);
                setError(null);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load builds');
            } finally {
                setLoading(false);
            }
        };
        loadBuilds();
    }, [projectId, suiteId]);

    const handleBuildClick = (buildId: string | number) => {
        setSelectedBuildId(buildId);
        onBuildSelect(buildId);
    };

    if (loading) {
        return <Spinner animation="border" />;
    }

    if (error) {
        return <Alert variant="danger">{error}</Alert>;
    }

    return (
        <div className="build-select-container">
            <div className="build-select">
                {builds.map((build) => (
                    <div
                        key={build.id}
                        className={`build-select-item ${build.id === selectedBuildId ? 'active' : ''}`}
                        onClick={() => handleBuildClick(build.id)}
                    >
                        {build.build_number}
                    </div>
                ))}
            </div>
        </div>
    );
};

export default BuildSelect;
