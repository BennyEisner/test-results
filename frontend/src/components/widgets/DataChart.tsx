import { useEffect, useRef, useCallback } from 'react';
import { Chart, registerables, ChartData, ChartOptions } from 'chart.js';
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

const colorPalette = [
    '#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#6366F1', '#EC4899', '#6B7280',
];
const borderPalette = [
    '#1D4ED8', '#059669', '#D97706', '#DC2626', '#7C3AED', '#4338CA', '#BE185D', '#4B5563',
];

const transformData = (data: DataChartDTO | null, chartType: 'line' | 'bar' | 'pie' | 'doughnut'): ChartData => {
    if (!data || !data.datasets) {
        return { labels: [], datasets: [] };
    }

    return {
        labels: data.labels || [],
        datasets: (data.datasets || []).map(dataset => {
            const isPieOrDoughnut = chartType === 'pie' || chartType === 'doughnut';
            const backgroundColors = isPieOrDoughnut
                ? (data.labels || []).map((_, i) => colorPalette[i % colorPalette.length])
                : dataset.backgroundColor || 'rgba(75, 192, 192, 0.6)';
            const borderColors = isPieOrDoughnut
                ? (data.labels || []).map((_, i) => borderPalette[i % borderPalette.length])
                : dataset.borderColor || 'rgba(75, 192, 192, 1)';

            return {
                ...dataset,
                backgroundColor: backgroundColors,
                borderColor: borderColors,
                borderWidth: 1,
            };
        }),
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
    console.log(`DataChart (${dataSource}): Props`, { projectId, suiteId, buildId, chartType, dataSource, isStatic, staticProjectId, staticSuiteId, staticBuildId, limit, refreshOn });
    const chartRef = useRef<HTMLCanvasElement>(null);
    const chartInstance = useRef<Chart | null>(null);

    const fetcher = useCallback((pid: number, sid?: number, bid?: number, lim?: number, signal?: AbortSignal) => {
        return dashboardApi.getChartData(pid, dataSource, sid, bid, lim, signal);
    }, [dataSource]);

    const { data, error, isLoading } = useSmartRefresh({
        projectId,
        suiteId,
        buildId,
        isStatic,
        staticProjectId,
        staticSuiteId,
        staticBuildId,
        fetcher,
        refreshOn,
        limit,
    });

    console.log(`DataChart (${dataSource}): render`, { isLoading, data, error });

    useEffect(() => {
        console.log(`DataChart (${dataSource}): useEffect`, { data });
        if (chartRef.current && data?.chart_data) {
            const transformedData = transformData(data.chart_data, chartType);
            console.log(`DataChart (${dataSource}): Transformed Chart Data:`, transformedData);

            if (chartInstance.current) {
                console.log(`DataChart (${dataSource}): Updating chart`);
                chartInstance.current.data = transformedData;
                chartInstance.current.options = getChartOptions(data.chart_data, chartType);
                chartInstance.current.update();
            } else {
                console.log(`DataChart (${dataSource}): Creating new chart`);
                chartInstance.current = new Chart(chartRef.current, {
                    type: chartType,
                    data: transformedData,
                    options: getChartOptions(data.chart_data, chartType),
                });
            }
        } else {
            console.log(`DataChart (${dataSource}): No data or chart ref, skipping chart update`);
        }

        return () => {
            if (chartInstance.current) {
                console.log(`DataChart (${dataSource}): Destroying chart`);
                chartInstance.current.destroy();
                chartInstance.current = null;
            }
        };
    }, [chartType, data, dataSource]);

    if (isLoading) {
        console.log(`DataChart (${dataSource}): Rendering Loading...`);
        return <div>Loading...</div>;
    }

    if (error) {
        console.log(`DataChart (${dataSource}): Rendering Error: ${error.message}`);
        return <div>Error: {error.message}</div>;
    }

    if (!data?.chart_data) {
        console.log(`DataChart (${dataSource}): Rendering No data`);
        return <div>No data available.</div>;
    }

    console.log(`DataChart (${dataSource}): Rendering canvas`);
    return <canvas ref={chartRef} />;
};

const getChartOptions = (data: DataChartDTO | null, chartType: 'line' | 'bar' | 'pie' | 'doughnut'): ChartOptions => {
    const isPieOrDoughnut = chartType === 'pie' || chartType === 'doughnut';

    if (isPieOrDoughnut) {
        return {
            scales: {
                x: { display: false },
                y: { display: false },
            },
            plugins: {
                legend: {
                    position: 'top' as const,
                },
                title: {
                    display: true,
                    text: data?.xAxisLabel || '',
                },
            },
        };
    }

    return {
        scales: {
            x: {
                title: {
                    display: !!data?.xAxisLabel,
                    text: data?.xAxisLabel || '',
                },
            },
            y: {
                title: {
                    display: !!data?.yAxisLabel,
                    text: data?.yAxisLabel || '',
                },
            },
        },
    };
};

export default DataChart;
