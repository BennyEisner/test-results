export interface DashboardLayout {
  id: string;
  name: string;
  components: DashboardComponents[];
  settings: DashboardSettings;
}

export interface DashboardComponents {
  id: string;
  type: ComponentType;
  positions: ComponentPosition;
  size: ComponentSize;
  props: ComponentProps;
  visible: boolean;
}
export interface ComponentPosition {
  x: number;
  y: number;
  order?: number;
}

export interface ComponentSize {
  width: number | "auto" | "full";
  height: number | "auto";
  minWidth?: number;
  minHeight?: number;
}

export interface ComponentProps {
  title?: string;
  filters?: Record<string, any>;
  dataSource?: string;
  chartType?: "doughnut" | "line" | "bar"; // Will add more once more visual components are added
  tableColumns?: string[];
  refreshInterval?: number;
  [key: string]: any;
}

export type ComponentType =
  | "builds-table"
  | "executions-summary"
  | "build-chart"
  | "test-cases-table"
  | "failures-summary"
  | "project-overview"
  | "search-bar"
  | "custom-widget";

export interface DashboardSettings {
  theme: "light" | "dark";
  layout: "grid" | "flex" | "masonry";
  spacing: "compact" | "normal" | "spacious";
  autoRefresh: boolean;
  refreshInterval: number;
}
export interface DashboardSettings {
  id: string;
  name: string;
  components: DashboardComponents[];
  settings: DashboardSettings;
}
