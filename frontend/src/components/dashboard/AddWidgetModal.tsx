import { useState } from 'react';
import { Modal, Nav } from 'react-bootstrap';
import { ComponentType, ComponentProps } from '../../types/dashboard';
import { COMPONENT_DEFINITIONS } from './ComponentRegistry';
import './AddWidgetModal.css';

interface AddWidgetModalProps {
    isOpen: boolean;
    onClose: () => void;
    onAddComponent: (type: ComponentType, props?: ComponentProps, isStatic?: boolean) => void;
}

const AddWidgetModal = ({ isOpen, onClose, onAddComponent }: AddWidgetModalProps) => {
    const [selectedTab, setSelectedTab] = useState<'static' | 'dynamic'>('dynamic');

    const staticComponents = Object.entries(COMPONENT_DEFINITIONS).filter(
        ([, def]) => def.configFields && def.configFields.length > 0
    );

    const handleSelect = (type: ComponentType, isStatic: boolean) => {
        const componentDef = COMPONENT_DEFINITIONS[type];
        onAddComponent(type, componentDef.defaultProps, isStatic);
        onClose();
    }

    return (
        <Modal show={isOpen} onHide={onClose} centered>
            <Modal.Header closeButton>
                <Modal.Title>Add Widget</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Nav variant="tabs" activeKey={selectedTab} onSelect={(k) => setSelectedTab(k as 'static' | 'dynamic')}>
                    <Nav.Item>
                        <Nav.Link eventKey="dynamic">Dynamic</Nav.Link>
                    </Nav.Item>
                    <Nav.Item>
                        <Nav.Link eventKey="static">Static</Nav.Link>
                    </Nav.Item>
                </Nav>
                <div className="mt-3">
                    {selectedTab === 'dynamic' && (
                        <div className="component-list">
                            {Object.entries(COMPONENT_DEFINITIONS).map(([type, def]) => (
                                <button
                                    key={type}
                                    className="widget-btn"
                                    onClick={() => handleSelect(type as ComponentType, false)}
                                >
                                    {def.name}
                                </button>
                            ))}
                        </div>
                    )}
                    {selectedTab === 'static' && (
                        <div className="component-list">
                            {staticComponents.map(([type, def]) => (
                                <button
                                    key={type}
                                    className="widget-btn"
                                    onClick={() => handleSelect(type as ComponentType, true)}
                                >
                                    {def.name}
                                </button>
                            ))}
                        </div>
                    )}
                </div>
            </Modal.Body>
        </Modal>
    );
};

export default AddWidgetModal;
