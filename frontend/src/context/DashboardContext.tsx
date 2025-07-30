import { createContext, useState, useContext, ReactNode } from 'react';

interface DashboardContextType {
    projectId: string | number | null;
    suiteId: string | number | null;
    buildId: string | number | null;
    setProjectId: (id: string | number | null) => void;
    setSuiteId: (id: string | number | null) => void;
    setBuildId: (id: string | number | null) => void;
}

const DashboardContext = createContext<DashboardContextType | undefined>(undefined);

export const DashboardProvider = ({ children }: { children: ReactNode }) => {
    const [projectId, setProjectId] = useState<string | number | null>(null);
    const [suiteId, setSuiteId] = useState<string | number | null>(null);
    const [buildId, setBuildId] = useState<string | number | null>(null);

    return (
        <DashboardContext.Provider
            value={{
                projectId,
                suiteId,
                buildId,
                setProjectId,
                setSuiteId,
                setBuildId,
            }}
        >
            {children}
        </DashboardContext.Provider>
    );
};

export const useDashboard = () => {
    const context = useContext(DashboardContext);
    if (context === undefined) {
        throw new Error('useDashboard must be used within a DashboardProvider');
    }
    return context;
};
