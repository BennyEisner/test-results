import { useState, useEffect } from 'react';
import { ListGroup, Spinner, Alert } from 'react-bootstrap';
import type { Suite } from '../types';
import { fetchSuites } from '../services/api';
import './SuiteMenu.css';

interface SuiteMenuProps {
    projectId: string | number;
    onSuiteSelect: (suiteId: string | number) => void;
}

const SuiteMenu = ({ projectId, onSuiteSelect }: SuiteMenuProps) => {
    const [suites, setSuites] = useState<Suite[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [selectedSuiteId, setSelectedSuiteId] = useState<string | number | null>(null);

    useEffect(() => {
        const loadSuites = async () => {
            try {
                setLoading(true);
                const data = await fetchSuites(projectId);
                setSuites(data);
                setError(null);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load suites');
            } finally {
                setLoading(false);
            }
        };
        loadSuites();
    }, [projectId]);

    const handleSuiteClick = (suiteId: string | number) => {
        setSelectedSuiteId(suiteId);
        onSuiteSelect(suiteId);
    };

    if (loading) {
        return <Spinner animation="border" />;
    }

    if (error) {
        return <Alert variant="danger">{error}</Alert>;
    }

    return (
        <div className="suite-menu">
            <h5 className="suite-menu-title">Test Suites</h5>
            <ListGroup>
                {suites.map((suite) => (
                    <ListGroup.Item
                        key={suite.id}
                        action
                        active={suite.id === selectedSuiteId}
                        onClick={() => handleSuiteClick(suite.id)}
                    >
                        {suite.name}
                    </ListGroup.Item>
                ))}
            </ListGroup>
        </div>
    );
};

export default SuiteMenu;
