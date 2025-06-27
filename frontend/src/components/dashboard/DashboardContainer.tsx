import { Responsive, WidthProvider } from 'react-grid-layout';
import { DashboardLayout, GridLayoutItem } from '../../types/dashboard';
import ComponentRegistry from './ComponentRegistry';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import './DashboardContainer.css';

const ResponsiveGridLayout = WidthProvider(Responsive);

interface DashboardContainerProps {
  layout: DashboardLayout;
  isEditing?: boolean;
  onLayoutChange?: (gridLayout: GridLayoutItem[]) => void;
  onRemoveComponent?: (componentId: string) => void;
  projectId?: string | number;
}

const DashboardContainer = ({
  layout,
  isEditing = false,
  onLayoutChange,
  onRemoveComponent,
  projectId,
}: DashboardContainerProps) => {
  const handleGridLayoutChange = (newGridLayout: any[]) => {
    if (onLayoutChange) {
      onLayoutChange(newGridLayout);
    }
  };

  const visibleComponents = layout.components.filter(comp => comp.visible);
  const visibleGridLayout = layout.gridLayout.filter(item => visibleComponents.some(c => c.id === item.i));

  return (
    <ResponsiveGridLayout
      className="dashboard-grid"
      layouts={{ lg: visibleGridLayout }}
      breakpoints={{ lg: 1200, md: 996, sm: 768, xs: 480, xxs: 0 }}
      cols={{ lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 }}
      rowHeight={60}
      isDraggable={isEditing}
      isResizable={isEditing}
      onLayoutChange={handleGridLayoutChange}
      margin={[16, 16]}
    >
      {visibleComponents.map((component) => (
        <div key={component.id} className="dashboard-item">
          {isEditing && (
            <div className="component-header">
              <span className="component-title">{component.props.title || component.type}</span>
              <button 
                className="remove-btn"
                onClick={() => onRemoveComponent && onRemoveComponent(component.id)}
              >
                x
              </button>
            </div>
          )}
          <div className="component-content">
            <ComponentRegistry type={component.type} props={component.props} projectId={projectId} />
          </div>
        </div>
      ))}
    </ResponsiveGridLayout>
  );
};

export default DashboardContainer;
