import axios from "axios";
import { DashboardLayout } from "../types/dashboard";

const API_BASE = "/api";
const USER_ID = 1; // Hardcoded user ID for develepemnt

export const dashboardApi = {
  getLayouts: async (): Promise<{
    layouts: DashboardLayout[];
    activeId: string;
  }> => {
    const response = await axios.get(`${API_BASE}/users/${USER_ID}/configs`);
    const layouts = JSON.parse(response.data.layouts);
    return { layouts, activeId: response.data.active_layout_id };
  },

  saveLayouts: async (
    layouts: DashboardLayout[],
    activeId: string,
  ): Promise<void> => {
    await axios.post(`${API_BASE}/users/${USER_ID}/configs`, {
      layouts,
      activeLayoutId: activeId,
    });
  },
};
