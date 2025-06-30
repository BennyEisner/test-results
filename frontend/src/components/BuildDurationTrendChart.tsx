import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { getBuildDurationTrends } from '../services/api';
import { BuildDurationTrend } from '../types';
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

const BuildDurationTrendChart: React.FC<BuildDurationTrendChartProps> = ({ projectId, suiteId }) => {
  const [chartData, setChartData] = useState<any>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const trends = await getBuildDurationTrends(projectId, suiteId);
        const data = {
          labels: trends.map((t: BuildDurationTrend) => new Date(t.created_at).toLocaleDateString()),
          datasets: [
            {
              label: 'Build Duration (s)',
              data: trends.map((t: BuildDurationTrend) => t.duration),
              fill: false,
              borderColor: 'rgb(75, 192, 192)',
              tension: 0.1,
            },
          ],
        };
        setChartData(data);
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
