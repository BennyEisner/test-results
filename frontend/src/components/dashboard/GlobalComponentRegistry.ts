import { ComponentDefinition } from "../../types/dashboard";
import { fetchRecentBuilds, fetchMostFailedTests } from "../../services/api";

export const GLOBAL_COMPONENT_DEFINITIONS: {
  [key: string]: ComponentDefinition;
} = {
  "builds-table": {
    name: "Builds Table",
    description: "Displays a table of recent builds.",
    category: "Tables",
    defaultProps: { title: "Recent Builds", fetchFunction: fetchRecentBuilds },
    defaultGridSize: { w: 4, h: 6 },
    configFields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        placeholder: "Enter a title",
      },
      { key: "projectId", label: "Project", type: "select" },
    ],
  },
  "build-duration-trend-chart": {
    name: "Build Duration Trend Chart",
    description: "Displays a trend chart of build durations.",
    category: "Charts",
    defaultProps: { title: "Build Duration Trend" },
    defaultGridSize: { w: 4, h: 6 },
    configFields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        placeholder: "Enter a title",
      },
      { key: "projectId", label: "Project", type: "select" },
      { key: "suiteId", label: "Suite", type: "select" },
    ],
  },
  "build-chart": {
    name: "Build Chart",
    description: "Displays a doughnut chart of build status.",
    category: "Charts",
    defaultProps: { title: "Build Status" },
    defaultGridSize: { w: 3, h: 6 },
    configFields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        placeholder: "Enter a title",
      },
      {
        key: "buildId",
        label: "Build ID",
        type: "text",
        placeholder: "Enter a build ID",
      },
    ],
  },
  "most-failed-tests-table": {
    name: "Most Failed Tests Table",
    description: "Displays a table of the most failed tests.",
    category: "Tables",
    defaultProps: {
      title: "Most Failed Tests",
      limit: 10,
      fetchMostFailedTests,
    },
    defaultGridSize: { w: 5, h: 5 },
    configFields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        placeholder: "Enter a title",
      },
      { key: "projectId", label: "Project", type: "select" },
      {
        key: "limit",
        label: "Limit",
        type: "text",
        placeholder: "Enter a limit",
      },
    ],
  },
  "executions-summary": {
    name: "Executions Summary",
    description: "Displays a summary of test executions.",
    category: "Summaries",
    defaultProps: { title: "Test Summary" },
    defaultGridSize: { w: 6, h: 3 },
    configFields: [
      {
        key: "title",
        label: "Title",
        type: "text",
        placeholder: "Enter a title",
      },
      {
        key: "buildId",
        label: "Build ID",
        type: "text",
        placeholder: "Enter a build ID",
      },
    ],
  },
};
