import { ComponentType } from '../../types/dashboard';
import { COMPONENT_DEFINITIONS } from './ComponentRegistry';
import './DashboardEditor.css';

interface DashboardEditorProps {
  onAddComponent: (type: ComponentType) => void;
}

const DashboardEditor = ({ onAddComponent }: DashboardEditorProps) => {
  return (
    <div className="dashboard-editor">
      <h4>Add Widget</h4>
      <div className="component-list">
        {Object.entries(COMPONENT_DEFINITIONS).map(([type, def]) => (
          <button 
            key={type} 
            className="add-component-btn"
            onClick={() => onAddComponent(type as ComponentType)}
          >
            + {def.name}
          </button>
        ))}
      </div>
    </div>
  );
};

export default DashboardEditor;
