import { useEffect, useRef, useState } from 'react';
import { Chart, registerables, ChartData } from 'chart.js';
import { dashboardApi } from '../../services/dashboardApi';
import { DataChartDTO } from '../../types/dashboard';

Chart.register(...registerables);

interface DataChartProps {
    projectId?: string | number;
    suiteId?: string | number;
    buildId?: string | number | null;
    chartType: 'line' | 'bar' | 'pie' | 'doughnut';
    dataSource: string;
    isStatic?: boolean;
    staticProjectId?: string | number;
    staticSuiteId?: string | number;
    staticBuildId?: string | number;
}

const transformData = (data: DataChartDTO | null): ChartData => {
    if (!data || !data.datasets) {
        return { labels: [], datasets: [] };
    }

    return {
        labels: data.labels,
        datasets: data.datasets.map(dataset => ({
            label: dataset.label,
            data: dataset.data,
            backgroundColor: 'rgba(75, 192, 192, 0.6)',
            borderColor: 'rgba(75, 192, 192, 1)',
            borderWidth: 1,
        })),
    };
};

const DataChart = ({
    projectId,
    suiteId,
    buildId,
    chartType,
    dataSource,
    isStatic,
    staticProjectId,
    staticSuiteId,
    staticBuildId,
}: DataChartProps) => {
    const chartRef = useRef<HTMLCanvasElement>(null);
    const chartInstance = useRef<Chart | null>(null);
    const [chartData, setChartData] = useState<DataChartDTO | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<Error | null>(null);

    const effectiveProjectId = isStatic ? staticProjectId : projectId;
    const effectiveSuiteId = isStatic ? staticSuiteId : suiteId;
    const effectiveBuildId = isStatic ? staticBuildId : buildId;

    useEffect(() => {
        if (effectiveProjectId) {
            const fetchChartData = async () => {
                try {
                    setLoading(true);
                    const response = await dashboardApi.getChartData(
                        Number(effectiveProjectId),
                        dataSource,
                        effectiveSuiteId ? Number(effectiveSuiteId) : undefined,
                        effectiveBuildId ? Number(effectiveBuildId) : undefined
                    );
                    setChartData(response.chart_data);
                } catch (err) {
                    setError(err as Error);
                } finally {
                    setLoading(false);
                }
            };
            fetchChartData();
        }
    }, [effectiveProjectId, effectiveSuiteId, effectiveBuildId, dataSource, isStatic]);

    useEffect(() => {
        if (chartRef.current) {
            const transformedData = transformData(chartData);
            console.log('Transformed Chart Data:', transformedData);

            if (chartInstance.current) {
                chartInstance.current.data = transformedData;
                chartInstance.current.update();
            } else {
                chartInstance.current = new Chart(chartRef.current, {
                    type: chartType,
                    data: transformedData,
                    options: {},
                });
            }
        }

        return () => {
            if (chartInstance.current) {
                chartInstance.current.destroy();
                chartInstance.current = null;
            }
        };
    }, [chartType, chartData]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return <canvas ref={chartRef} />;
};

export default DataChart;
