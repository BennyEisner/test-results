import { useState } from 'react';
import { DashboardContext } from '../../context/DashboardContext';
import SuiteMenu from '../suite/SuiteMenu';
import DashboardContainer from '../dashboard/DashboardContainer';
import DashboardEditor from '../dashboard/DashboardEditor';
import GlobalDashboardEditor from '../dashboard/GlobalDashboardEditor';
import { useDashboardLayouts } from '../../hooks/useDashboardLayouts';
import PageLayout from '../common/PageLayout';
import './HomePage.css';
import { DashboardLayout, ComponentType } from '../../types/dashboard';

const DashboardPage = () => {
    const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
    const [selectedSuiteId, setSelectedSuiteId] = useState<number | null>(null);
    const {
        activeLayout,
        isEditing,
        setIsEditing,
        updateGridLayout,
        addComponent: addAutoComponent,
        removeComponent,
        updateLayout,
    } = useDashboardLayouts();

    // Separate state for the global layout
    const [globalLayout, setGlobalLayout] = useState<DashboardLayout>({
        id: 'global',
        name: 'Global',
        components: [],
        gridLayout: [],
        settings: { theme: 'light', layout: 'grid', spacing: 'normal' },
    });

    const [editingMode, setEditingMode] = useState<'auto' | 'global' | null>(null);

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

    const addGlobalComponent = (type: ComponentType, props: any) => {
        const newId = `${type}-${Date.now()}`;
        const newComponent = {
            id: newId,
            type,
            props,
            visible: true,
        };
        const newLayoutItem = {
            i: newId,
            x: 0,
            y: 0,
            w: 4,
            h: 6,
        };
        setGlobalLayout(prevLayout => ({
            ...prevLayout,
            components: [...prevLayout.components, newComponent],
            gridLayout: [...prevLayout.gridLayout, newLayoutItem],
        }));
    };

    const removeGlobalComponent = (widgetId: string) => {
        setGlobalLayout(prevLayout => ({
            ...prevLayout,
            components: prevLayout.components.filter(c => c.id !== widgetId),
            gridLayout: prevLayout.gridLayout.filter(item => item.i !== widgetId),
        }));
    };

    const updateGlobalGridLayout = (gridLayout: any) => {
        setGlobalLayout(prevLayout => ({
            ...prevLayout,
            gridLayout,
        }));
    };

    const handleEdit = (mode: 'auto' | 'global') => {
        setIsEditing(true);
        setEditingMode(mode);
    }

    const handleDone = () => {
        setIsEditing(false);
        setEditingMode(null);
    }


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
                        {isEditing ? (
                            <button onClick={handleDone}>Done</button>
                        ) : (
                            <>
                                <button onClick={() => handleEdit('auto')}>Edit Project Dashboard</button>
                                <button onClick={() => handleEdit('global')}>Edit Global Widgets</button>
                            </>
                        )}
                    </div>

                    {isEditing && editingMode === 'auto' && (
                        <DashboardEditor onAddComponent={addAutoComponent} />
                    )}
                    {isEditing && editingMode === 'global' && (
                        <GlobalDashboardEditor onAddComponent={addGlobalComponent} />
                    )}

                    <h3>Project Dashboard</h3>
                    <DashboardContainer
                        layout={activeLayout}
                        isEditing={isEditing && editingMode === 'auto'}
                        onLayoutChange={updateGridLayout}
                        onRemoveComponent={removeComponent}
                        projectId={selectedProjectId ?? undefined}
                        suiteId={selectedSuiteId ?? undefined}
                    />

                    <h3>Global Widgets</h3>
                    <DashboardContainer
                        layout={globalLayout}
                        isEditing={isEditing && editingMode === 'global'}
                        onLayoutChange={updateGlobalGridLayout}
                        onRemoveComponent={removeGlobalComponent}
                    />
                </div>
            </PageLayout>
        </DashboardContext.Provider>
    );
};

export default DashboardPage;
