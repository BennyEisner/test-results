import axios from "axios";
import { DashboardLayout } from "../types/dashboard";

const API_BASE = "/api";
const USER_ID = 1; // Hardcoded user ID for development

export const dashboardApi = {
  getLayouts: async (): Promise<{
    layouts: DashboardLayout[];
    activeId: string;
  }> => {
    console.log("dashboardApi.getLayouts: Making API call to get layouts");
    const response = await axios.get(
      `${API_BASE}/users/${USER_ID}/configs`,
    );
    console.log("dashboardApi.getLayouts: Raw API response:", response.data);
    
    if (response.data && response.data.length > 0) {
      const config = response.data[0];
      console.log("dashboardApi.getLayouts: Found config:", config);
      const layouts = JSON.parse(config.layouts);
      console.log("dashboardApi.getLayouts: Parsed layouts:", layouts);
      const result = { layouts, activeId: config.active_layout_id };
      console.log("dashboardApi.getLayouts: Returning:", result);
      return result;
    }
    console.log("dashboardApi.getLayouts: No config found, returning empty");
    return { layouts: [], activeId: "" };
  },

  saveLayouts: async (
    layouts: DashboardLayout[],
    activeId: string,
  ): Promise<{
    layouts: DashboardLayout[];
    activeId: string;
  }> => {
    console.log("dashboardApi.saveLayouts: Saving layouts:", { layouts, activeId });
    const payload = {
      layouts: JSON.stringify(layouts),
      active_layout_id: activeId,
    };
    console.log("dashboardApi.saveLayouts: Payload:", payload);
    
    const response = await axios.post(`${API_BASE}/users/${USER_ID}/configs`, payload);
    console.log("dashboardApi.saveLayouts: Raw API response:", response.data);
    
    if (response.data) {
      const savedLayouts = JSON.parse(response.data.layouts);
      const result = { 
        layouts: savedLayouts, 
        activeId: response.data.active_layout_id 
      };
      console.log("dashboardApi.saveLayouts: Returning:", result);
      return result;
    }
    
    // Fallback to the input data if response doesn't contain expected data
    console.log("dashboardApi.saveLayouts: No response data, falling back to input");
    return { layouts, activeId };
  },

  saveActiveLayoutId: async (activeId: string): Promise<void> => {
    console.log("dashboardApi.saveActiveLayoutId: Saving active layout ID:", activeId);
    await axios.put(`${API_BASE}/users/${USER_ID}/configs/active`, {
      active_layout_id: activeId,
    });
    console.log("dashboardApi.saveActiveLayoutId: Successfully saved active layout ID");
  },
};
