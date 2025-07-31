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

const getChartColors = () => {
    const style = getComputedStyle(document.documentElement);
    return [
        style.getPropertyValue('--chart-color-1').trim(),
        style.getPropertyValue('--chart-color-2').trim(),
        style.getPropertyValue('--chart-color-3').trim(),
        style.getPropertyValue('--chart-color-4').trim(),
        style.getPropertyValue('--chart-color-5').trim(),
    ];
};

const transformData = (data: DataChartDTO | null, chartType: 'line' | 'bar' | 'pie' | 'doughnut'): ChartData => {
    if (!data || !data.datasets) {
        return { labels: [], datasets: [] };
    }

    const chartColors = getChartColors();

    return {
        labels: data.labels || [],
        datasets: (data.datasets || []).map(dataset => {
            const isPieOrDoughnut = chartType === 'pie' || chartType === 'doughnut';
            
            // Use dynamic colors from API if available, otherwise fall back to default
            const backgroundColors = Array.isArray(dataset.backgroundColor) && dataset.backgroundColor.length > 0
                ? dataset.backgroundColor
                : isPieOrDoughnut
                    ? (data.labels || []).map((_, i) => chartColors[i % chartColors.length])
                    : 'rgba(139, 233, 253, 0.6)';

            const borderColors = Array.isArray(dataset.borderColor) && dataset.borderColor.length > 0
                ? dataset.borderColor
                : isPieOrDoughnut
                    ? (data.labels || []).map((_, i) => chartColors[i % chartColors.length])
                    : 'rgba(139, 233, 253, 1)';

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
    const style = getComputedStyle(document.documentElement);
    const textColor = style.getPropertyValue('--text-primary').trim();
    const gridColor = style.getPropertyValue('--border-color').trim();

    const baseOptions: ChartOptions = {
        maintainAspectRatio: false,
        plugins: {
            legend: {
                position: 'top' as const,
                labels: {
                    color: textColor,
                },
            },
            tooltip: {
                backgroundColor: style.getPropertyValue('--background-secondary').trim(),
                titleColor: textColor,
                bodyColor: textColor,
                borderColor: gridColor,
                borderWidth: 1,
            },
        },
    };

    if (chartType === 'pie' || chartType === 'doughnut') {
        return {
            ...baseOptions,
            scales: {
                x: { display: false },
                y: { display: false },
            },
            plugins: {
                ...baseOptions.plugins,
                title: {
                    display: true,
                    text: data?.xAxisLabel || '',
                    color: textColor,
                },
            },
        };
    }

    return {
        ...baseOptions,
        scales: {
            x: {
                title: {
                    display: !!data?.xAxisLabel,
                    text: data?.xAxisLabel || '',
                    color: textColor,
                },
                ticks: {
                    color: textColor,
                },
                grid: {
                    color: gridColor,
                },
            },
            y: {
                title: {
                    display: !!data?.yAxisLabel,
                    text: data?.yAxisLabel || '',
                    color: textColor,
                },
                ticks: {
                    color: textColor,
                },
                grid: {
                    color: gridColor,
                },
            },
        },
    };
};

export default DataChart;
