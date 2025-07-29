import { Line } from 'react-chartjs-2';
import { useState, useEffect } from 'react';
import { dashboardApi } from '../../services/dashboardApi';
import { DataChartDTO } from '../../types/dashboard';
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js';

ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

interface BuildDurationTrendChartProps {
    projectId: number;
    suiteId: number;
}

const BuildDurationTrendChart = ({ projectId, suiteId }: BuildDurationTrendChartProps) => {
    const [chartData, setChartData] = useState<any>(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await dashboardApi.getChartData(projectId, 'build-duration-trend', suiteId);
                const trends = response.chart_data as DataChartDTO;
                if (trends && trends.labels && trends.datasets && trends.datasets.length > 0) {
                    setChartData(trends);
                } else {
                    console.log('No trends data available or invalid format:', trends);
                }
            } catch (error) {
                console.error('Error fetching build duration trends:', error);
            }
        };

        if (projectId && suiteId) {
            fetchData();
        }
    }, [projectId, suiteId]);

    return (
        <div>
            <h2>Build Duration Trend</h2>
            {chartData ? <Line data={chartData} /> : <p>Loading chart...</p>}
        </div>
    );
};

export default BuildDurationTrendChart;
