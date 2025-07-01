import { Doughnut } from 'react-chartjs-2';
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from 'chart.js';
import { useExecutionsSummary } from '../hooks/useExecutionsSummary';
import { Spinner } from 'react-bootstrap';

ChartJS.register(ArcElement, Tooltip, Legend);

interface BuildDoughnutChartProps {
    buildId?: string | number;
}

const BuildDoughnutChart = ({ buildId }: BuildDoughnutChartProps) => {
    const { stats, loading } = useExecutionsSummary(buildId);
    if (!buildId) {
        return <p className="text-center text-muted">No build selected.</p>;
    }
    if (loading) {
        return (
            <div className="d-flex justify-content-center align-items-center" style={{ height: '200px' }}>
                <Spinner animation="border" role="status" variant="primary">
                    <span className="visually-hidden">Loading...</span>
                </Spinner>
            </div>
        );
    }

    if (!stats || stats.total === 0) {
        return <p className="text-center text-muted">No data available to display chart.</p>;
    }

    const chartData = {
        labels: ['Passed', 'Failed', 'Skipped'],
        datasets: [
            {
                label: 'Test Executions',
                data: [stats.passed, stats.failed, stats.skipped],
                backgroundColor: [
                    'rgba(40, 167, 69, 0.8)',  // Success (green)
                    'rgba(220, 53, 69, 0.8)',   // Fail (red)
                    'rgba(108, 117, 125, 0.8)'  // Skipped (gray)
                ],
                borderColor: [
                    'rgba(40, 167, 69, 1)',
                    'rgba(220, 53, 69, 1)',
                    'rgba(108, 117, 125, 1)'
                ],
                borderWidth: 1,
            },
        ],
    };

    const chartOptions = {
        responsive: true,
        plugins: {
            legend: {
                position: 'top' as const,
            },
            title: {
                display: true,
                text: 'Executions Summary',
                font: {
                    size: 16
                }
            },
        },
    };

    return (
        <div style={{ height: '300px', display: 'flex', justifyContent: 'center' }}>
            <Doughnut data={chartData} options={chartOptions} />
        </div>
    );
};

export default BuildDoughnutChart;
