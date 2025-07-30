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
    const [asyncOptions, setAsyncOptions] = useState<Record<string, { value: string; label: string }[]>>({});

    const [projects, setProjects] = useState<{ id: number; name: string }[]>([]);

    useEffect(() => {
        if (isOpen && componentType) {
            const componentDef = COMPONENT_DEFINITIONS[componentType];
            const fieldsToFetch = componentDef.configFields?.filter(f => f.asyncOptions) || [];

            fieldsToFetch.forEach(field => {
                field.asyncOptions!().then(options => {
                    setAsyncOptions(prev => ({ ...prev, [field.key]: options }));
                });
            });

            fetchProjects().then(setProjects);
            setConfig(prev => ({ ...prev, isStatic: initialIsStatic }));
        }
    }, [isOpen, componentType, initialIsStatic]);

    if (!componentType) {
        return null;
    }

    const componentDef = COMPONENT_DEFINITIONS[componentType];

    const handleSave = () => {
        const finalProps = { ...componentDef.defaultProps, ...config };
        onSave(componentType, finalProps, config.isStatic || false);
        onClose();
    };

    const renderField = (field: ConfigField) => {
        if (field.condition && !field.condition(config)) {
            return null;
        }

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
                            {field.key === 'projectId' && <option value="">Select a Project</option>}
                            {field.key === 'projectId' ? projects.map((p) => (
                                <option key={p.id} value={p.id}>
                                    {p.name}
                                </option>
                            )) : (field.asyncOptions ? asyncOptions[field.key] : field.options)?.map(opt => {
                                if (typeof opt === 'object') {
                                    return <option key={opt.value} value={opt.value}>{opt.label}</option>;
                                }
                                return <option key={opt} value={opt}>{String(opt)}</option>;
                            })}
                        </Form.Select>
                        {helpText && <Form.Text className="text-muted">{helpText}</Form.Text>}
                    </Form.Group>
                );
            case 'textarea':
                return (
                    <Form.Group key={key} className="mb-3">
                        <Form.Label>{label}</Form.Label>
                        <Form.Control
                            as="textarea"
                            rows={3}
                            value={value}
                            placeholder={placeholder}
                            onChange={(e) => setConfig({ ...config, [key]: e.target.value })}
                        />
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
            case 'number':
                return (
                    <Form.Group key={key} className="mb-3">
                        <Form.Label>{label}</Form.Label>
                        <Form.Control
                            type="number"
                            value={value}
                            placeholder={placeholder}
                            onChange={(e) => setConfig({ ...config, [key]: e.target.value ? Number(e.target.value) : undefined })}
                        />
                    </Form.Group>
                );
            case 'checkbox':
                return (
                    <Form.Group key={key} className="mb-3">
                        <Form.Check
                            type="switch"
                            id={key}
                            label={label}
                            checked={value}
                            onChange={(e) => setConfig({ ...config, [key]: e.target.checked })}
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
                    {componentDef.configFields?.map(renderField)}
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
