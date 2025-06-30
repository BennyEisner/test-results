import { useState } from 'react';
import { DashboardContext } from '../context/DashboardContext';
import SuiteMenu from './SuiteMenu';
import DashboardContainer from './dashboard/DashboardContainer';
import DashboardEditor from './dashboard/DashboardEditor';
import { useDashboardLayouts } from '../hooks/useDashboardLayouts';
import PageLayout from './PageLayout';
import './HomePage.css';

const DashboardPage = () => {
  const {
    activeLayout,
    isEditing,
    setIsEditing,
    updateGridLayout,
    addComponent,
    removeComponent,
    updateLayout,
  } = useDashboardLayouts();
  const [selectedProjectId, setSelectedProjectId] = useState<number | null>(null);
  const [selectedSuiteId, setSelectedSuiteId] = useState<number | null>(null);

  const handleProjectSelect = (projectId: number) => {
    setSelectedSuiteId(null);
    setSelectedProjectId(projectId);
    if (activeLayout) {
      const updatedComponents = activeLayout.components.map(c => {
        const newProps = { ...c.props, projectId: projectId };
        if (c.props && 'suiteId' in c.props) {
          newProps.suiteId = undefined;
        }
        return { ...c, props: newProps };
      });
      updateLayout({ ...activeLayout, components: updatedComponents });
    }
  };

  const handleSuiteSelect = (suiteId: number | string) => {
    setSelectedSuiteId(typeof suiteId === 'string' ? parseInt(suiteId, 10) : suiteId);
    if (activeLayout) {
      const updatedComponents = activeLayout.components.map(c => {
        if (c.props && 'suiteId' in c.props) {
          return { ...c, props: { ...c.props, suiteId: suiteId } };
        }
        return c;
      });
      updateLayout({ ...activeLayout, components: updatedComponents });
    }
  };

  if (!activeLayout) {
    return <div>Loading dashboard...</div>;
  }

  return (
    <DashboardContext.Provider value={{ selectedProjectId, selectedSuiteId, onProjectSelect: handleProjectSelect, onSuiteSelect: handleSuiteSelect }}>
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
        <DashboardEditor onAddComponent={addComponent} />
      )}

      <DashboardContainer
        layout={activeLayout}
        isEditing={isEditing}
        onLayoutChange={updateGridLayout}
        onRemoveComponent={removeComponent}
        projectId={selectedProjectId ?? undefined}
      />
        </div>
      </PageLayout>
    </DashboardContext.Provider>
  );
};

export default DashboardPage;
