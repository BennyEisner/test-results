import { useState } from 'react';
import { ComponentType, ComponentProps } from '../../types/dashboard';
import { GLOBAL_COMPONENT_DEFINITIONS } from './GlobalComponentRegistry';
import ComponentConfigModal from './ComponentConfigModal';
import './DashboardEditor.css';

interface GlobalDashboardEditorProps {
  onAddComponent: (type: ComponentType, props?: ComponentProps) => void;
}

const GlobalDashboardEditor = ({ onAddComponent }: GlobalDashboardEditorProps) => {
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [selectedComponentType, setSelectedComponentType] = useState<ComponentType | null>(null);

  const handleComponentSelect = (componentType: ComponentType) => {
    const componentDef = GLOBAL_COMPONENT_DEFINITIONS[componentType];
    
    if (componentDef.configFields && componentDef.configFields.length > 0) {
      // Show configuration modal
      setSelectedComponentType(componentType);
      setShowConfigModal(true);
    } else {
      // Add component directly with default props
      onAddComponent(componentType, componentDef.defaultProps);
    }
  };

  const handleConfigSave = (componentType: ComponentType, props: ComponentProps) => {
    onAddComponent(componentType, props);
    setShowConfigModal(false);
    setSelectedComponentType(null);
  };

  const handleConfigClose = () => {
    setShowConfigModal(false);
    setSelectedComponentType(null);
  };

  return (
    <div className="dashboard-editor">
      <h4>Add Widget</h4>
      <div className="component-list">
        {Object.entries(GLOBAL_COMPONENT_DEFINITIONS).map(([type, def]) => (
          <button 
            key={type} 
            className="add-component-btn"
            onClick={() => handleComponentSelect(type as ComponentType)}
          >
            + {def.name}
          </button>
        ))}
      </div>
      
      <ComponentConfigModal
        isOpen={showConfigModal}
        onClose={handleConfigClose}
        componentType={selectedComponentType}
        onSave={handleConfigSave}
      />
    </div>
  );
};

export default GlobalDashboardEditor;
