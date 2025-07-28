import api from './api';
import { StatusBadgeDTO, MetricCardDTO, DataChartDTO, DashboardLayout } from '../types/dashboard';

interface ChartDataResponse {
    chart_data: DataChartDTO;
}

export const getStatus = async (projectId: number): Promise<StatusBadgeDTO> => {
  const response = await api.get(`/dashboard/projects/${projectId}/status`);
  return response.data;
};

export const getMetric = async (projectId: number, metricType: string): Promise<MetricCardDTO> => {
  const response = await api.get(`/dashboard/projects/${projectId}/metric/${metricType}`);
  return response.data;
};

export const getChartData = async (
    projectId: number,
    chartType: string,
    suiteId?: number,
    buildId?: number,
    limit?: number,
): Promise<ChartDataResponse> => {
    const params = new URLSearchParams();
    if (suiteId) {
        params.append('suite_id', String(suiteId));
    }
    if (buildId) {
        params.append('build_id', String(buildId));
    }
    if (limit) {
        params.append('limit', String(limit));
    }
    const queryString = params.toString();
    const url = `/dashboard/projects/${projectId}/chart/${chartType}${queryString ? `?${queryString}` : ''}`;
    const response = await api.get(url);
    return response.data;
};

export const getAvailableWidgets = async (): Promise<any> => {
    const response = await api.get('/dashboard/available-widgets');
    return response.data;
};

export const getLayouts = async (userId: number) => {
    const response = await api.get(`/users/${userId}/config`);
    return response.data;
};

export const saveLayouts = async (userId: number, layouts: DashboardLayout[], activeId: string) => {
    const response = await api.put(`/users/${userId}/config`, { 
        layouts: JSON.stringify(layouts), 
        active_layout_id: activeId 
    });
    return response.data;
};

export const saveActiveLayoutId = async (userId: number, activeId: string) => {
    const response = await api.put(`/users/${userId}/config`, { 
        active_layout_id: activeId 
    });
    return response.data;
};

export const dashboardApi = {
    getStatus,
    getMetric,
    getChartData,
    getAvailableWidgets,
    getLayouts,
    saveLayouts,
    saveActiveLayoutId,
};
