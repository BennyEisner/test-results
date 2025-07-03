import { useState } from 'react';
import { DashboardContext } from '../../context/DashboardContext';
import SuiteMenu from '../suite/SuiteMenu';
import DashboardContainer from '../dashboard/DashboardContainer';
import DashboardEditor from '../dashboard/DashboardEditor';
import { useDashboardLayouts } from '../../hooks/useDashboardLayouts';
import PageLayout from '../common/PageLayout';
import './HomePage.css';
import { ComponentType, ComponentProps } from '../../types/dashboard';

const DashboardPage = () => {
    const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
    const [selectedSuiteId, setSelectedSuiteId] = useState<number | null>(null);
    const {
        activeLayout,
        isEditing,
        setIsEditing,
        updateGridLayout,
        addComponent,
        removeComponent,
        updateLayout,
    } = useDashboardLayouts();

    const handleProjectSelect = (projectId: number) => {
        setSelectedProjectId(projectId);
        setSelectedSuiteId(null);
    };

    const handleSuiteSelect = (suiteId: number | string) => {
        setSelectedSuiteId(typeof suiteId === 'string' ? parseInt(suiteId, 10) : suiteId);
    };

    const updateWidgetProps = (widgetId: string, props: Record<string, any>) => {
        if (activeLayout) {
            const updatedComponents = activeLayout.components.map(c => {
                if (c.id === widgetId) {
                    return { ...c, props: { ...c.props, ...props } };
                }
                return c;
            });
            updateLayout({ ...activeLayout, components: updatedComponents });
        }
    };

    const handleAddComponent = (type: ComponentType, props?: ComponentProps, isStatic?: boolean) => {
        addComponent(type, props, isStatic);
    };

    if (!activeLayout) {
        return <div>Loading dashboard...</div>;
    }

    return (
        <DashboardContext.Provider value={{ selectedProjectId, selectedSuiteId, onProjectSelect: handleProjectSelect, onSuiteSelect: handleSuiteSelect, updateWidgetProps }}>
            <PageLayout>
                {selectedProjectId && (
                    <SuiteMenu
                        projectId={selectedProjectId}
                        onSuiteSelect={handleSuiteSelect}
                    />
                )}
                <div className={`home-page ${selectedProjectId ? 'dashboard-with-sidebar' : ''}`}>
                    <div className="dashboard-header">
                        <h2>{activeLayout.name}{selectedProjectId && ` - Project ${selectedProjectId}`}</h2>
                        <button onClick={() => setIsEditing(!isEditing)}>
                            {isEditing ? 'Done' : 'Edit Dashboard'}
                        </button>
                    </div>

                    {isEditing && (
                        <DashboardEditor onAddComponent={handleAddComponent} />
                    )}

                    <DashboardContainer
                        layout={activeLayout}
                        isEditing={isEditing}
                        onLayoutChange={updateGridLayout}
                        onRemoveComponent={removeComponent}
                        projectId={selectedProjectId ?? undefined}
                        suiteId={selectedSuiteId ?? undefined}
                    />
                </div>
            </PageLayout>
        </DashboardContext.Provider>
    );
};

export default DashboardPage;
