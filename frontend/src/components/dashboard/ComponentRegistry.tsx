import { ComponentType, ComponentProps, ComponentDefinition } from '../../types/dashboard';
import BuildsTable from '../BuildsTable';
import ExecutionsSummary from '../ExecutionsSummary';
import BuildDoughnutChart from '../BuildDoughnutChart';
import { fetchRecentBuilds } from '../../services/api';

interface ComponentRegistryProps {
  type: ComponentType;
  props: ComponentProps;
}

const ComponentRegistry = ({ type, props }: ComponentRegistryProps) => {
  const componentProps = { ...props, className: `dashboard-component--${type}` };

  switch (type) {
    case 'builds-table':
      return <BuildsTable {...componentProps} />;
    case 'executions-summary':
      return <ExecutionsSummary {...componentProps} />;
    case 'build-chart':
      return <BuildDoughnutChart {...componentProps} />;
    default:
      return <div className="component-placeholder">Unknown component: {type}</div>;
  }
};

export default ComponentRegistry;

export const COMPONENT_DEFINITIONS: Record<ComponentType, ComponentDefinition> = {
  'builds-table': {
    name: 'Builds Table',
    description: 'Display recent builds in a table format',
    category: 'Tables',
    defaultProps: { title: 'Recent Builds', fetchFunction: fetchRecentBuilds },
    defaultGridSize: { w: 8, h: 6, minW: 4, minH: 4 },
  },
  'executions-summary': {
    name: 'Executions Summary',
    description: 'Summary statistics of test executions',
    category: 'Summaries',
    defaultProps: { title: 'Test Executions' },
    defaultGridSize: { w: 4, h: 4, minW: 3, minH: 3 },
  },
  'build-chart': {
    name: 'Build Chart',
    description: 'Visual chart of build data',
    category: 'Charts',
    defaultProps: { title: 'Build Status' },
    defaultGridSize: { w: 6, h: 5, minW: 4, minH: 4 },
  },
};
