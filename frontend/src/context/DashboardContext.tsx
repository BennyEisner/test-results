import { createContext, useContext } from 'react';

interface DashboardContextType {
  selectedProjectId: number | null;
  selectedSuiteId: number | null;
  onProjectSelect: (projectId: number) => void;
  onSuiteSelect: (suiteId: number | string) => void;
}

export const DashboardContext = createContext<DashboardContextType | undefined>(undefined);

export const useDashboard = () => {
  return useContext(DashboardContext);
};
