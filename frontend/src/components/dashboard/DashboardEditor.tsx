import { useState } from 'react';
import { ComponentType, ComponentProps } from '../../types/dashboard';
import { COMPONENT_DEFINITIONS } from './ComponentRegistry';
import ComponentConfigModal from './ComponentConfigModal';
import AddWidgetModal from './AddWidgetModal';
import './DashboardEditor.css';

interface DashboardEditorProps {
  onAddComponent: (type: ComponentType, props?: ComponentProps, isStatic?: boolean) => void;
}

const DashboardEditor = ({ onAddComponent }: DashboardEditorProps) => {
  const [showAddModal, setShowAddModal] = useState(false);
  const [showConfigModal, setShowConfigModal] = useState(false);
  const [selectedComponentType, setSelectedComponentType] = useState<ComponentType | null>(null);
  const [isStatic, setIsStatic] = useState(false);

  const handleAddComponent = (type: ComponentType, props?: ComponentProps, isStatic?: boolean) => {
    const componentDef = COMPONENT_DEFINITIONS[type];
    if (isStatic && componentDef.configFields && componentDef.configFields.length > 0) {
      setSelectedComponentType(type);
      setIsStatic(true);
      setShowConfigModal(true);
    } else {
      onAddComponent(type, props, isStatic);
    }
    setShowAddModal(false);
  };

  const handleConfigSave = (componentType: ComponentType, props: ComponentProps, isStatic: boolean) => {
    onAddComponent(componentType, props, isStatic);
    setShowConfigModal(false);
    setSelectedComponentType(null);
  };

  const handleConfigClose = () => {
    setShowConfigModal(false);
    setSelectedComponentType(null);
    setIsStatic(false);
  };

  return (
    <div className="dashboard-editor">
      <button className="btn btn-primary add-widget-button" onClick={() => setShowAddModal(true)}>
        Add Widget
      </button>
      
      <AddWidgetModal
        isOpen={showAddModal}
        onClose={() => setShowAddModal(false)}
        onAddComponent={handleAddComponent}
      />

      <ComponentConfigModal
        isOpen={showConfigModal}
        onClose={handleConfigClose}
        componentType={selectedComponentType}
        onSave={handleConfigSave}
        initialIsStatic={isStatic}
      />
    </div>
  );
};

export default DashboardEditor;
