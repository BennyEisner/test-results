import { useState, useEffect } from "react";
import {
  DashboardLayout,
  GridLayoutItem,
  ComponentType,
  ComponentProps,
} from "../types/dashboard";
import { COMPONENT_DEFINITIONS } from "../components/dashboard/ComponentRegistry";
import { dashboardApi } from "../services/dashboardApi";
import { useAuth } from "../context/AuthContext";

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
      props: { title: "Recent Builds" },
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
  const { user } = useAuth();
  const [layouts, setLayouts] = useState<DashboardLayout[]>([defaultLayout]);
  const [activeLayoutId, setActiveLayoutId] = useState<string>("default");
  const [isEditing, setIsEditing] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isInitialized, setIsInitialized] = useState(false);

  // Load layouts from API or localStorage
  useEffect(() => {
    const loadLayouts = async () => {
      setIsLoading(true);
      setError(null);
      console.log("Attempting to load layouts...");

      if (user) {
        try {
          // Try API first
          console.log("Fetching layouts from API...");
          const data = await dashboardApi.getLayouts(user.id);
          console.log("API response received:", data);

          if (data) {
            let finalLayouts: DashboardLayout[];
            let finalActiveId: string;

            if (data.layouts) {
              try {
                finalLayouts = JSON.parse(data.layouts);
                if (!Array.isArray(finalLayouts) || finalLayouts.length === 0) {
                  finalLayouts = [defaultLayout];
                }
              } catch (e) {
                console.error("Failed to parse layouts from API, using default.", e);
                finalLayouts = [defaultLayout];
              }
            } else {
              finalLayouts = [defaultLayout];
            }

            finalActiveId = data.active_layout_id || finalLayouts[0]?.id || "default";

            console.log("Setting final layouts:", finalLayouts);
            setLayouts(finalLayouts);
            setActiveLayoutId(finalActiveId);

            // Update localStorage with server data
            console.log("Updating localStorage with server data:", {
              layouts: finalLayouts,
              activeLayoutId: finalActiveId,
            });
            localStorage.setItem(
              STORAGE_KEY,
              JSON.stringify({
                layouts: finalLayouts,
                activeLayoutId: finalActiveId,
              }),
            );

            setIsLoading(false);
            setIsInitialized(true);
            console.log("Initialization complete - API data loaded");
            return;
          }
        } catch (e) {
          console.error(
            "Failed to load layouts from API, using localStorage fallback",
            e,
          );
          setError("Could not connect to the server. Displaying cached data.");
          // Continue to localStorage fallback
        }
      }

      // Fallback to localStorage
      console.log("Falling back to localStorage.");
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        console.log("Found stored layouts in localStorage:", stored);
        try {
          const parsed = JSON.parse(stored);
          if (
            parsed.layouts &&
            parsed.layouts.length > 0 &&
            parsed.layouts[0].version >= LAYOUT_VERSION
          ) {
            console.log("Applying layouts from localStorage.");
            setLayouts(parsed.layouts);
            setActiveLayoutId(parsed.activeLayoutId || parsed.layouts[0].id);
          } else {
            console.log(
              "localStorage data is outdated or empty, using default layout.",
            );
            setLayouts([defaultLayout]);
            setActiveLayoutId("default");
          }
        } catch (e) {
          console.error(
            "Failed to parse layouts from localStorage, using default.",
            e,
          );
          localStorage.removeItem(STORAGE_KEY);
          setLayouts([defaultLayout]);
          setActiveLayoutId("default");
        }
      } else {
        console.log("No layouts in localStorage, using default.");
        setLayouts([defaultLayout]);
        setActiveLayoutId("default");
      }
      setIsLoading(false);
      setIsInitialized(true);
      console.log("Initialization complete - fallback data loaded");
    };

    loadLayouts();
  }, [user]);

  const saveLayouts = async (
    newLayouts: DashboardLayout[],
    newActiveId: string,
  ) => {
    if (!isInitialized) {
      console.warn(
        "Attempted to save layouts before initialization is complete. Aborting.",
      );
      return;
    }

    if (!user) {
      console.warn("Attempted to save layouts without a user. Aborting.");
      return;
    }
    const previousLayouts = layouts;
    const previousActiveId = activeLayoutId;

    console.log("Attempting to save layouts:", { newLayouts, newActiveId });
    // Optimistically update the state
    setLayouts(newLayouts);
    setActiveLayoutId(newActiveId);
    setIsSaving(true);
    setError(null);

    try {
      console.log("Sending layouts to API...");
      const savedData = await dashboardApi.saveLayouts(
        user.id,
        newLayouts,
        newActiveId,
      );
      console.log("API save response received:", savedData);
      if (savedData) {
        let savedLayouts: DashboardLayout[];
        try {
          savedLayouts = JSON.parse(savedData.layouts);
        } catch (e) {
          console.error("Failed to parse saved layouts, using optimistic state.", e);
          savedLayouts = newLayouts; 
        }

        const savedActiveId = savedData.active_layout_id || newActiveId;

        // Update state with the response from the server
        console.log("Updating state and localStorage with server response.");
        setLayouts(savedLayouts);
        setActiveLayoutId(savedActiveId);
        localStorage.setItem(
          STORAGE_KEY,
          JSON.stringify({
            layouts: savedLayouts,
            activeLayoutId: savedActiveId,
          }),
        );
      }
    } catch (e) {
      console.error("Failed to save layouts to API, rolling back.", e);
      setError(
        "Failed to save changes to the server. Your changes have been reverted.",
      );

      // Rollback to the previous state
      console.log("Rolling back state and localStorage.");
      setLayouts(previousLayouts);
      setActiveLayoutId(previousActiveId);

      // Also rollback localStorage to ensure consistency
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          layouts: previousLayouts,
          activeLayoutId: previousActiveId,
        }),
      );
    } finally {
      setIsSaving(false);
    }
  };

  const changeActiveLayoutId = async (layoutId: string) => {
    if (!isInitialized || !user) {
      console.warn(
        "Attempted to change active layout before initialization is complete or without a user. Aborting.",
      );
      return;
    }

    const previousActiveId = activeLayoutId;
    setActiveLayoutId(layoutId);

    try {
      await dashboardApi.saveActiveLayoutId(user.id, layoutId);
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          layouts: layouts,
          activeLayoutId: layoutId,
        }),
      );
    } catch (e) {
      console.error("Failed to save active layout ID, rolling back.", e);
      setError("Failed to save active layout. Your change has been reverted.");
      setActiveLayoutId(previousActiveId);
    }
  };

  const updateLayout = (updatedLayout: DashboardLayout) => {
    const newLayouts = layouts.map((l) =>
      l.id === updatedLayout.id ? updatedLayout : l,
    );
    saveLayouts(newLayouts, activeLayoutId);
  };
  const updateGridLayout = (gridLayout: GridLayoutItem[]) => {
    if (!isInitialized) {
      return;
    }
    const activeLayout = layouts.find((l) => l.id === activeLayoutId);
    if (activeLayout) {
      const updatedLayout = { ...activeLayout, gridLayout };
      const newLayouts = layouts.map((l) =>
        l.id === updatedLayout.id ? updatedLayout : l,
      );
      saveLayouts(newLayouts, activeLayoutId);
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
