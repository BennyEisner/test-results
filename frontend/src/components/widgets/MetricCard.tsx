import { useState, useEffect } from 'react';
import { dashboardApi } from '../../services/dashboardApi';
import { MetricCardDTO } from '../../types/dashboard';
import './MetricCard.css';

interface MetricCardProps {
    title?: string;
    projectId?: string | number;
    metricType: string;
}

const MetricCard = ({ title, projectId, metricType }: MetricCardProps) => {
    const [metric, setMetric] = useState<MetricCardDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (projectId) {
            const fetchMetric = async () => {
                try {
                    setLoading(true);
                    const data = await dashboardApi.getMetric(Number(projectId), metricType);
                    setMetric(data);
                } catch (err) {
                    setError(err as Error);
                } finally {
                    setLoading(false);
                }
            };
            fetchMetric();
        }
    }, [projectId, metricType]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return (
        <div className="metric-card">
            <h3 className="metric-title">{title || metric?.metric}</h3>
            <div className="metric-value">{metric?.value}</div>
        </div>
    );
};

export default MetricCard;
