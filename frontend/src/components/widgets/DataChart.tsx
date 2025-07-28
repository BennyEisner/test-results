import { useEffect, useRef } from 'react';
import { Chart, registerables, ChartData } from 'chart.js';
import { dashboardApi } from '../../services/dashboardApi';
import { DataChartDTO } from '../../types/dashboard';
import { useSmartRefresh } from '../../hooks/useSmartRefresh';

Chart.register(...registerables);

type RefreshTrigger = 'project' | 'suite' | 'build';

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
    limit?: number;
    refreshOn?: RefreshTrigger[];
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
    limit,
    refreshOn = ['project', 'suite', 'build'],
}: DataChartProps) => {
    const chartRef = useRef<HTMLCanvasElement>(null);
    const chartInstance = useRef<Chart | null>(null);

    const { data, error, isLoading } = useSmartRefresh({
        projectId,
        suiteId,
        buildId,
        isStatic,
        staticProjectId,
        staticSuiteId,
        staticBuildId,
        fetcher: (pid, sid, bid, lim) => dashboardApi.getChartData(pid, dataSource, sid, bid, lim),
        refreshOn,
        limit,
    });

    useEffect(() => {
        if (chartRef.current && data) {
            const transformedData = transformData(data.chart_data);
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
    }, [chartType, data]);

    if (isLoading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error.message}</div>;
    }

    return <canvas ref={chartRef} />;
};

export default DataChart;
