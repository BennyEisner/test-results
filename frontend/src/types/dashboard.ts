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
  fetchFunction?: () => Promise<any[]>;
  onResultSelect?: (result: any) => void;
  [key: string]: any;
}

export interface ConfigField {
  key: string;
  label: string;
  type: "select" | "text" | "number" | "boolean" | "multi-select";
  required?: boolean;
  options?: { value: string | number; label: string }[];
  defaultValue?: any;
  placeholder?: string;
  helpText?: string;
}

export type ComponentType =
  | "builds-table"
  | "build-chart"
  | "build-duration-trend-chart"
  | "most-failed-tests-table"
  | "executions-summary";

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
