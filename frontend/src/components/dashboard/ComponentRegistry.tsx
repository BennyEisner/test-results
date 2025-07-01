import { ComponentType, ComponentProps, ComponentDefinition } from '../../types/dashboard';
import BuildsTable from '../BuildsTable';
import ExecutionsSummary from '../ExecutionsSummary';
import BuildDoughnutChart from '../BuildDoughnutChart';
import BuildDurationTrendChart from '../BuildDurationTrendChart';
import MostFailedTestsTable from '../MostFailedTestsTable';
import { fetchRecentBuilds } from '../../services/api';

interface ComponentRegistryProps {
    type: ComponentType;
    props: ComponentProps;
    projectId?: string | number;
}

function ComponentRegistry({ type, props, projectId }: ComponentRegistryProps) {
    const componentProps = { ...props, className: `dashboard-component--${type}` };

    // Handle project-specific configurations
    if (type === 'builds-table') {
        const fetchProjectId = props.projectId && props.projectId !== 'all' ? props.projectId : projectId;
        if (fetchProjectId) {
            componentProps.fetchFunction = () => fetchRecentBuilds(fetchProjectId);
        } else {
            componentProps.fetchFunction = () => fetchRecentBuilds();
        }
    }

    switch (type) {
        case 'builds-table':
            return <BuildsTable {...componentProps} />;
        case 'build-chart':
            return <BuildDoughnutChart {...componentProps} />;
        case 'build-duration-trend-chart':
            const finalProjectId = props.projectId ? Number(props.projectId) : (projectId ? Number(projectId) : undefined);
            const suiteIdNumber = props.suiteId ? Number(props.suiteId) : undefined;
            if (finalProjectId && suiteIdNumber) {
                const { projectId: _, suiteId: __, ...restProps } = componentProps;
                return <BuildDurationTrendChart projectId={finalProjectId} suiteId={suiteIdNumber} {...restProps} />;
            }
            return <div className="component-placeholder">Project ID and Suite ID required for Build Duration Trend Chart</div>;

        case 'most-failed-tests-table':
            const finalMostFailedProjectId = props.projectId ? Number(props.projectId) : (projectId ? Number(projectId) : undefined);
            const mostFailedSuiteIdNumber = props.suiteId ? Number(props.suiteId) : undefined;
            const limit = props.limit ? Number(props.limit) : 10;
            if (finalMostFailedProjectId) {
                const { projectId: _, suiteId: __, ...restProps } = componentProps;
                return <MostFailedTestsTable projectId={finalMostFailedProjectId} limit={limit} suiteId={mostFailedSuiteIdNumber} {...restProps} />;
            }
            return <div className="component-placeholder">Project ID required for Most Failed Tests Table</div>;
        case 'executions-summary':
            return <ExecutionsSummary {...componentProps} />;

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
                required: true,
                placeholder: 'Enter build ID'
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
                required: true,
                placeholder: 'Enter project ID'
            },
            {
                key: 'suiteId',
                label: 'Suite ID',
                type: 'text',
                required: true,
                placeholder: 'Enter suite ID'
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
                required: true,
                placeholder: 'Enter project ID'
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
                required: true,
                placeholder: 'Enter build ID'
            }
        ]
    },

};
