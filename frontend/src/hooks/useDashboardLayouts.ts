import { useState, useEffect } from "react";
import {
  DashboardLayout,
  GridLayoutItem,
  ComponentType,
  ComponentProps,
} from "../types/dashboard";
import { COMPONENT_DEFINITIONS } from "../components/dashboard/ComponentRegistry";
import { fetchRecentBuilds } from "../services/api";
import { dashboardApi } from "../services/dashboardApi";

const STORAGE_KEY = "dashboard-layouts";

const LAYOUT_VERSION = 2;

const defaultLayout: DashboardLayout = {
  id: "default",
  name: "Default Dashboard",
  version: LAYOUT_VERSION,
  components: [
    {
      id: "builds-1",
      type: "builds-table",
      props: { title: "Recent Builds", fetchFunction: fetchRecentBuilds },
      visible: true,
    },
    {
      id: "build-duration-trend-chart-1",
      type: "build-duration-trend-chart",
      props: { title: "Build Duration Trend", projectId: 1, suiteId: 1 },
      visible: true,
    },

    {
      id: "chart-1",
      type: "build-chart",
      props: { title: "Build Status", buildId: "1" },
      visible: true,
    },
    {
      id: "most-failed-tests-table-1",
      type: "most-failed-tests-table",
      props: { title: "Most Failed Tests", projectId: 1, limit: 10 },
      visible: true,
    },
    {
      id: "summary-1",
      type: "executions-summary",
      props: { title: "Test Summary", buildId: "1" },
      visible: true,
    },
  ],
  gridLayout: [
    { i: "builds-1", x: 0, y: 0, w: 4, h: 6 },
    { i: "chart-1", x: 8, y: 0, w: 3, h: 6 },
    { i: "build-duration-trend-chart-1", x: 4, y: 0, w: 4, h: 6 },
    { i: "most-failed-tests-table-1", x: 0, y: 6, w: 5, h: 5 },
    { i: "summary-1", x: 8, y: 6, w: 6, h: 3 },
  ],
  settings: { theme: "light", layout: "grid", spacing: "normal" },
};

export const useDashboardLayouts = () => {
  const [layouts, setLayouts] = useState<DashboardLayout[]>([defaultLayout]);
  const [activeLayoutId, setActiveLayoutId] = useState<string>("default");
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load layouts from API or localStorage
  useEffect(() => {
    const loadLayouts = async () => {
      setIsLoading(true);
      setError(null);
      try {
        // Try API first
        const data = await dashboardApi.getLayouts();
        if (data && data.layouts && data.layouts.length > 0) {
          setLayouts(data.layouts);
          setActiveLayoutId(data.activeId || data.layouts[0].id);

          // Update localStorage with server data
          localStorage.setItem(
            STORAGE_KEY,
            JSON.stringify({
              layouts: data.layouts,
              activeLayoutId: data.activeId || data.layouts[0].id,
            }),
          );

          setIsLoading(false);
          return;
        }
      } catch (e) {
        console.warn(
          "Failed to load layouts from API, using localStorage fallback",
          e,
        );
        setError("Could not connect to the server. Displaying cached data.");
        // Continue to localStorage fallback
      }

      // Fallback to localStorage
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        try {
          const parsed = JSON.parse(stored);
          if (
            parsed.layouts &&
            parsed.layouts.length > 0 &&
            parsed.layouts[0].version >= LAYOUT_VERSION
          ) {
            setLayouts(parsed.layouts);
            setActiveLayoutId(parsed.activeLayoutId || parsed.layouts[0].id);
          } else {
            // If localStorage is empty or outdated, use default
            setLayouts([defaultLayout]);
            setActiveLayoutId("default");
          }
        } catch (e) {
          console.error("Failed to parse layouts from localStorage", e);
          localStorage.removeItem(STORAGE_KEY);
          setLayouts([defaultLayout]);
          setActiveLayoutId("default");
        }
      } else {
        // If no localStorage, use default
        setLayouts([defaultLayout]);
        setActiveLayoutId("default");
      }
      setIsLoading(false);
    };

    loadLayouts();
  }, []);

  const saveLayouts = async (
    newLayouts: DashboardLayout[],
    newActiveId: string,
  ) => {
    // Update state immediately
    setLayouts(newLayouts);
    setActiveLayoutId(newActiveId);

    // Always save to localStorage
    localStorage.setItem(
      STORAGE_KEY,
      JSON.stringify({ layouts: newLayouts, activeLayoutId: newActiveId }),
    );

    // Try API
    setIsSaving(true);
    setError(null);
    try {
      await dashboardApi.saveLayouts(newLayouts, newActiveId);
    } catch (e) {
      console.error("Failed to save layouts to API", e);
      setError("Changes saved locally but couldn't sync with server.");
    } finally {
      setIsSaving(false);
    }
  };

  const changeActiveLayoutId = (layoutId: string) => {
    saveLayouts(layouts, layoutId);
  };

  const updateLayout = (updatedLayout: DashboardLayout) => {
    const newLayouts = layouts.map((l) =>
      l.id === updatedLayout.id ? updatedLayout : l,
    );
    saveLayouts(newLayouts, activeLayoutId);
  };
  const updateGridLayout = (gridLayout: GridLayoutItem[]) => {
    const activeLayout = layouts.find((l) => l.id === activeLayoutId);
    if (activeLayout) {
      updateLayout({ ...activeLayout, gridLayout });
    }
  };

  const addComponent = (
    type: ComponentType,
    props?: ComponentProps,
    isStatic?: boolean,
  ) => {


    const activeLayout = layouts.find((l) => l.id === activeLayoutId);
    if (!activeLayout) return;

    const definition = COMPONENT_DEFINITIONS[type];
    const newId = `${type}-${Date.now()}`;

    const newComponent = {
      id: newId,
      type,
      props: props || definition.defaultProps,
      visible: true,
      isStatic: isStatic || false,
    };

    const newLayoutItem = {
      i: newId,
      x: 0,
      y: 0, // Add to top
      ...definition.defaultGridSize,
    };

    // Move all other items down
    const updatedGridLayout = activeLayout.gridLayout.map((item) => ({
      ...item,
      y: item.y + (definition.defaultGridSize.h || 1),
    }));

    updateLayout({
      ...activeLayout,
      components: [...activeLayout.components, newComponent],
      gridLayout: [...updatedGridLayout, newLayoutItem],
    });
  };

  const removeComponent = (componentId: string) => {
    const activeLayout = layouts.find((l) => l.id === activeLayoutId);
    if (!activeLayout) return;

    updateLayout({
      ...activeLayout,
      components: activeLayout.components.filter((c) => c.id !== componentId),
      gridLayout: activeLayout.gridLayout.filter(
        (item) => item.i !== componentId,
      ),
    });
  };

  const activeLayout =
    layouts.find((l) => l.id === activeLayoutId) || layouts[0];

  return {
    layouts,
    activeLayout,
    updateLayout,
    saveLayouts,
    setActiveLayoutId: changeActiveLayoutId,
    isEditing,
    setIsEditing,
    isLoading,
    isSaving,
    error,
    updateGridLayout,
    addComponent,
    removeComponent,
  };
};
