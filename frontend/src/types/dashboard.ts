export interface DashboardLayout {
  id: string;
  name: string;
  version?: number;
  components: DashboardComponent[];
  gridLayout: GridLayoutItem[];
  settings: DashboardSettings;
}

export interface DashboardComponent {
  id: string;
  type: ComponentType;
  props: ComponentProps;
  visible: boolean;
  isStatic?: boolean;
}

// matches react-grid-layout's expected format
export interface GridLayoutItem {
  i: string; // component id
  x: number;
  y: number;
  w: number;
  h: number;
  minW?: number;
  minH?: number;
  maxW?: number;
  maxH?: number;
  static?: boolean; // prevents dragging/resizing
}

export interface ComponentProps {
  title?: string;
  buildId?: string | number;
  projectId?: string | number;
  suiteId?: string | number;
  dataSource?: string;
  fetchFunction?: () => Promise<any[]>;
  onResultSelect?: (result: any) => void;
  [key: string]: any;
}

export interface ConfigField {
    key: string;
    label: string;
    type: 'text' | 'select' | 'number' | 'checkbox' | 'textarea';
    options?: string[] | { value: string; label: string }[];
    asyncOptions?: () => Promise<{ value: string; label: string }[]>;
    required?: boolean;
    defaultValue?: any;
    placeholder?: string;
    helpText?: string;
    condition?: (props: ComponentProps) => boolean;
}

export type ComponentType =
  | "builds-table"
  | "build-chart"
  | "build-duration-trend-chart"
  | "most-failed-tests-table"
  | "executions-summary"
  | "metric-card"
  | "status-badge"
  | "data-chart";

export interface DashboardSettings {
  theme: "light" | "dark";
  layout: "grid" | "flex";
  spacing: "compact" | "normal" | "spacious";
}

export interface ComponentDefinition {
  name: string;
  description: string;
  category: string;
  defaultProps: ComponentProps;
  defaultGridSize: {
    w: number;
    h: number;
    minW?: number;
    minH?: number;
    maxW?: number;
    maxH?: number;
  };
  configFields?: ConfigField[];
}

export interface StatusBadgeDTO {
  status: string;
  count: number;
}

export interface MetricCardDTO {
  metric: string;
  value: number;
}

export interface DataChartDTO {
  labels: string[];
  datasets: DatasetDTO[];
  xAxisLabel?: string;
  yAxisLabel?: string;
}

export interface DatasetDTO {
  label: string;
  data: number[];
  backgroundColor?: string | string[];
  borderColor?: string | string[];
}

export interface Widget {
    id: string;
    name: string;
    description: string;
    component: string;
    props: any;
}
