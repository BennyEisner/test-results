import DashboardContainer from './dashboard/DashboardContainer';
import DashboardEditor from './dashboard/DashboardEditor';
import { useDashboardLayouts } from '../hooks/useDashboardLayouts';
import './HomePage.css';

const DashboardPage = () => {
  const { 
    activeLayout, 
    isEditing, 
    setIsEditing,
    updateGridLayout,
    addComponent,
    removeComponent,
  } = useDashboardLayouts();

  if (!activeLayout) {
    return <div>Loading dashboard...</div>;
  }

  return (
    <div className="home-page">
      <div className="dashboard-header">
        <h2>{activeLayout.name}</h2>
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
      />
    </div>
  );
};

export default DashboardPage;
