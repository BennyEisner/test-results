import { useState, useEffect } from 'react';
import { Modal, Button, Form } from 'react-bootstrap';
import { ComponentType, ComponentProps, ConfigField } from '../../types/dashboard';
import { COMPONENT_DEFINITIONS } from './ComponentRegistry';
import { fetchProjects } from '../../services/api';

interface ComponentConfigModalProps {
    isOpen: boolean;
    onClose: () => void;
    componentType: ComponentType | null;
    onSave: (componentType: ComponentType, props: ComponentProps, isStatic: boolean) => void;
    initialIsStatic?: boolean;
}

const ComponentConfigModal = ({ isOpen, onClose, componentType, onSave, initialIsStatic = false }: ComponentConfigModalProps) => {
    const [config, setConfig] = useState<ComponentProps>({});
    const [isStatic, setIsStatic] = useState(initialIsStatic);

    const [projects, setProjects] = useState<{ id: number; name: string }[]>([]);

    useEffect(() => {
        if (isOpen) {
            fetchProjects().then(setProjects);
            setIsStatic(initialIsStatic);
        }
    }, [isOpen, initialIsStatic]);

    if (!componentType) {
        return null;
    }

    const componentDef = COMPONENT_DEFINITIONS[componentType];

    const handleSave = () => {
        const finalProps = { ...componentDef.defaultProps, ...config };
        onSave(componentType, finalProps, isStatic);
        onClose();
    };

    const renderField = (field: ConfigField) => {

        const { key, label, type, placeholder, helpText } = field;
        const value = config[key] || field.defaultValue;

        switch (type) {
            case 'select':
                return (
                    <Form.Group key={key} className="mb-3">
                        <Form.Label>{label}</Form.Label>
                        <Form.Select
                            value={value}
                            onChange={(e) => setConfig({ ...config, [key]: e.target.value })}
                        >
                            {field.key === 'projectId' ? (
                                <>
                                    <option value="">Select a Project</option>
                                    {projects.map((p) => (
                                        <option key={p.id} value={p.id}>
                                            {p.name}
                                        </option>
                                    ))}
                                </>
                            ) : (field.options ? (
                                field.options.map(opt => <option key={opt.value} value={opt.value}>{opt.label}</option>)
                            ) : null)}
                        </Form.Select>
                        {helpText && <Form.Text className="text-muted">{helpText}</Form.Text>}
                    </Form.Group>
                );
            case 'text':
                return (
                    <Form.Group key={key} className="mb-3">
                        <Form.Label>{label}</Form.Label>
                        <Form.Control
                            type="text"
                            value={value}
                            placeholder={placeholder}
                            onChange={(e) => setConfig({ ...config, [key]: e.target.value })}
                        />
                    </Form.Group>
                );
            default:
                return null;
        }
    };

    return (
        <Modal show={isOpen} onHide={onClose}>
            <Modal.Header closeButton>
                <Modal.Title>Configure {componentDef.name}</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form>
                    {isStatic && componentDef.configFields?.map(renderField)}
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={onClose}>
                    Cancel
                </Button>
                <Button variant="primary" onClick={handleSave}>
                    Save
                </Button>
            </Modal.Footer>
        </Modal>
    );
};

export default ComponentConfigModal;
