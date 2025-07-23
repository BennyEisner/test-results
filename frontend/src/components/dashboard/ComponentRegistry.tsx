import { ComponentType, ComponentProps, ComponentDefinition } from '../../types/dashboard';
import BuildsTable from '../build/BuildsTable';
import ExecutionsSummary from '../execution/ExecutionsSummary';
import BuildDoughnutChart from '../build/BuildDoughnutChart';
import BuildDurationTrendChart from '../build/BuildDurationTrendChart';
import MostFailedTestsTable from '../test/MostFailedTestsTable';
import MetricCard from '../widgets/MetricCard';
import StatusBadge from '../widgets/StatusBadge';
import DataChart from '../widgets/DataChart';
import { fetchBuilds } from '../../services/api';

interface ComponentRegistryProps {
    type: ComponentType;
    props: ComponentProps;
    projectId?: string | number;
    suiteId?: string | number;
    buildId?: string | number | null;
}

function ComponentRegistry({ type, props, projectId, suiteId, buildId }: ComponentRegistryProps) {
    const componentProps = { ...props, className: `dashboard-component--${type}` };

    // Handle project-specific configurations
    if (type === 'builds-table') {
        const fetchProjectId = props.projectId && props.projectId !== 'all' ? props.projectId : projectId;
        if (fetchProjectId) {
            componentProps.fetchFunction = () => fetchBuilds(fetchProjectId);
        } else {
            // Fallback: show a placeholder or use a default projectId (e.g., 1)
            componentProps.fetchFunction = () => Promise.resolve([]); // or fetchBuilds(1)
        }
    }

    switch (type) {
        case 'builds-table':
            return <BuildsTable {...componentProps} />;
        case 'build-chart':
            const chartBuildId = (buildId ? Number(buildId) : undefined) ?? (props.buildId ? Number(props.buildId) : undefined);
            if (chartBuildId) {
                return <BuildDoughnutChart buildId={chartBuildId} {...componentProps} />;
            }
            return <div className="component-placeholder">Select a build to view the chart.</div>;

        case 'build-duration-trend-chart':
            const finalProjectId = (projectId ? Number(projectId) : undefined) ?? (props.projectId ? Number(props.projectId) : undefined);
            const suiteIdNumber = (suiteId ? Number(suiteId) : undefined) ?? (props.suiteId ? Number(props.suiteId) : undefined);
            if (finalProjectId && suiteIdNumber) {
                const { projectId: _, suiteId: __, ...restProps } = componentProps;
                return <BuildDurationTrendChart projectId={finalProjectId} suiteId={suiteIdNumber} {...restProps} />;
            }
            return <div className="component-placeholder">Select a project and suite to view the trend chart.</div>;

        case 'most-failed-tests-table':
            const finalMostFailedProjectId = (projectId ? Number(projectId) : undefined) ?? (props.projectId ? Number(props.projectId) : undefined);
            const mostFailedSuiteIdNumber = (suiteId ? Number(suiteId) : undefined) ?? (props.suiteId ? Number(props.suiteId) : undefined);
            const limit = props.limit ? Number(props.limit) : 10;
            if (finalMostFailedProjectId) {
                const { projectId: _, suiteId: __, ...restProps } = componentProps;
                return <MostFailedTestsTable projectId={finalMostFailedProjectId} limit={limit} suiteId={mostFailedSuiteIdNumber} {...restProps} />;
            }
            return <div className="component-placeholder">Select a project to view the most failed tests.</div>;
        case 'executions-summary':
            const summaryBuildId = (buildId ? Number(buildId) : undefined) ?? (props.buildId ? Number(props.buildId) : undefined);
            if (summaryBuildId) {
                return <ExecutionsSummary buildId={summaryBuildId} {...componentProps} />;
            }
            return <div className="component-placeholder">Select a build to view the summary.</div>;
        
        case 'metric-card':
            return <MetricCard {...componentProps} value={props.value} />;

        case 'status-badge':
            return <StatusBadge {...componentProps} status={props.status} >{props.children}</StatusBadge>;

        case 'data-chart':
            return <DataChart {...componentProps} type={props.type} data={props.data} />;

        default:
            return <div className="component-placeholder">Unknown component: {type}</div>;
    }
}

export default ComponentRegistry;

export const COMPONENT_DEFINITIONS: Record<ComponentType, ComponentDefinition> = {
    'builds-table': {
        name: 'Builds Table',
        description: 'Display recent builds in a table format',
        category: 'Tables',
        defaultProps: { title: 'Recent Builds' },
        defaultGridSize: { w: 3, h: 6, minW: 4, minH: 4 },
        configFields: [
            {
                key: 'projectId',
                label: 'Project',
                type: 'select',
                required: false,
                defaultValue: 'all',
                helpText: 'Select a specific project or show all projects'
            },
            {
                key: 'title',
                label: 'Component Title',
                type: 'text',
                defaultValue: 'Recent Builds',
                placeholder: 'Enter component title'
            }
        ]
    },
    'build-chart': {
        name: 'Build Chart',
        description: 'Visual chart of build data',
        category: 'Charts',
        defaultProps: { title: 'Build Status' },
        defaultGridSize: { w: 5, h: 6, minW: 4, minH: 4 },
        configFields: [
            {
                key: 'title',
                label: 'Component Title',
                type: 'text',
                defaultValue: 'Build Status',
                placeholder: 'Enter component title'
            },
            {
                key: 'buildId',
                label: 'Build ID',
                type: 'text',
                required: false,
                placeholder: 'Enter build ID (optional)'
            }
        ]
    },
    'build-duration-trend-chart': {
        name: 'Build Duration Trend',
        description: 'Shows the trend of build durations over time.',
        category: 'Charts',
        defaultProps: { title: 'Build Duration Trend' },
        defaultGridSize: { w: 5, h: 7, minW: 4, minH: 4 },
        configFields: [
            {
                key: 'title',
                label: 'Component Title',
                type: 'text',
                defaultValue: 'Build Duration Trend',
                placeholder: 'Enter component title'
            },
            {
                key: 'projectId',
                label: 'Project ID',
                type: 'text',
                required: false,
                placeholder: 'Enter project ID (optional)'
            },
            {
                key: 'suiteId',
                label: 'Suite ID',
                type: 'text',
                required: false,
                placeholder: 'Enter suite ID (optional)'
            }
        ]
    },
    'most-failed-tests-table': {
        name: 'Most Failed Tests Table',
        description: 'Display the most frequently failing tests',
        category: 'Tables',
        defaultProps: { title: 'Most Failed Tests' },
        defaultGridSize: { w: 5, h: 6, minW: 4, minH: 4 },
        configFields: [
            {
                key: 'title',
                label: 'Component Title',
                type: 'text',
                defaultValue: 'Most Failed Tests',
                placeholder: 'Enter component title'
            },
            {
                key: 'projectId',
                label: 'Project ID',
                type: 'text',
                required: false,
                placeholder: 'Enter project ID (optional)'
            },
            {
                key: 'limit',
                label: 'Limit',
                type: 'number',
                required: false,
                defaultValue: 10,
                placeholder: 'Enter the number of tests to display'
            },
            {
                key: 'suiteId',
                label: 'Suite ID',
                type: 'text',
                required: false,
                placeholder: 'Enter suite ID (optional)'
            }
        ]
    },
    'executions-summary': {
        name: 'Executions Summary',
        description: 'Summary statistics of test executions',
        category: 'Summaries',
        defaultProps: { title: 'Test Executions' },
        defaultGridSize: { w: 4, h: 4, minW: 3, minH: 3 },
        configFields: [
            {
                key: 'title',
                label: 'Component Title',
                type: 'text',
                defaultValue: 'Test Executions',
                placeholder: 'Enter component title'
            },
            {
                key: 'buildId',
                label: 'Build ID',
                type: 'text',
                required: false,
                placeholder: 'Enter build ID (optional)'
            }
        ]
    },
    'metric-card': {
        name: 'Metric Card',
        description: 'Display a single metric with a trend indicator',
        category: 'Widgets',
        defaultProps: { title: 'Metric', value: 'N/A' },
        defaultGridSize: { w: 3, h: 2, minW: 2, minH: 2 },
        configFields: [
            {
                key: 'title',
                label: 'Title',
                type: 'text',
                defaultValue: 'Metric',
            },
            {
                key: 'value',
                label: 'Value',
                type: 'text',
                defaultValue: 'N/A',
            },
            {
                key: 'change',
                label: 'Change',
                type: 'text',
            },
            {
                key: 'changeType',
                label: 'Change Type',
                type: 'select',
                options: ['increase', 'decrease', 'neutral'],
                defaultValue: 'neutral',
            },
        ],
    },
    'status-badge': {
        name: 'Status Badge',
        description: 'Display a status badge',
        category: 'Widgets',
        defaultProps: { status: 'neutral', children: 'Status' },
        defaultGridSize: { w: 2, h: 1, minW: 2, minH: 1 },
        configFields: [
            {
                key: 'status',
                label: 'Status',
                type: 'select',
                options: ['success', 'warning', 'danger', 'info', 'neutral'],
                defaultValue: 'neutral',
            },
            {
                key: 'children',
                label: 'Text',
                type: 'text',
                defaultValue: 'Status',
            },
        ],
    },
    'data-chart': {
        name: 'Data Chart',
        description: 'Display a chart from data',
        category: 'Charts',
        defaultProps: { type: 'bar', data: {} },
        defaultGridSize: { w: 6, h: 8, minW: 4, minH: 4 },
        configFields: [
            {
                key: 'type',
                label: 'Chart Type',
                type: 'select',
                options: ['line', 'bar', 'pie', 'doughnut'],
                defaultValue: 'bar',
            },
            {
                key: 'data',
                label: 'Data',
                type: 'textarea',
                defaultValue: '{}',
            },
            {
                key: 'options',
                label: 'Options',
                type: 'textarea',
                defaultValue: '{}',
            },
        ],
    },
};
