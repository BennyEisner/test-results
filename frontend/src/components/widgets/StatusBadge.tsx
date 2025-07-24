import { useState, useEffect } from 'react';
import { dashboardApi } from '../../services/dashboardApi';
import { StatusBadgeDTO } from '../../types/dashboard';
import './StatusBadge.css';

interface StatusBadgeProps {
    projectId?: string | number;
}

const StatusBadge = ({ projectId }: StatusBadgeProps) => {
    const [status, setStatus] = useState<StatusBadgeDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    useEffect(() => {
        if (projectId) {
            const fetchStatus = async () => {
                try {
                    setLoading(true);
                    const data = await dashboardApi.getStatus(Number(projectId));
                    setStatus(data);
                } catch (err) {
                    setError(err as Error);
                } finally {
                    setLoading(false);
                }
            };
            fetchStatus();
        }
    }, [projectId]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    const getStatusClass = (status: string | undefined) => {
        switch (status) {
            case 'success':
                return 'status-badge-success';
            case 'warning':
                return 'status-badge-warning';
            case 'danger':
                return 'status-badge-danger';
            case 'info':
                return 'status-badge-info';
            default:
                return 'status-badge-neutral';
        }
    };

    return (
        <span className={`status-badge ${getStatusClass(status?.status)}`}>
            {status?.status}
        </span>
    );
};

export default StatusBadge;
