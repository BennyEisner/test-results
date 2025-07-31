import { useDashboard } from '../../context/DashboardContext';
import SuiteMenu from '../suite/SuiteMenu';
import DashboardContainer from '../dashboard/DashboardContainer';
import DashboardEditor from '../dashboard/DashboardEditor';
import { useDashboardLayouts } from '../../hooks/useDashboardLayouts';
import PageLayout from '../common/PageLayout';
import './HomePage.css';
import { ComponentType, ComponentProps } from '../../types/dashboard';

const DashboardPage = () => {
    const { projectId, setProjectId, setSuiteId, setBuildId } = useDashboard();
    const {
        activeLayout,
        isEditing,
        setIsEditing,
        updateGridLayout,
        addComponent,
        removeComponent,
    } = useDashboardLayouts();

    const handleProjectSelect = (projectId: number) => {
        setProjectId(projectId);
        setSuiteId(null);
        setBuildId(null);
    };

    const handleSuiteSelect = (suiteId: number | string) => {
        setSuiteId(typeof suiteId === 'string' ? parseInt(suiteId, 10) : suiteId);
        setBuildId(null);
    };

    const handleAddComponent = (type: ComponentType, props?: ComponentProps, isStatic?: boolean) => {
        addComponent(type, props, isStatic);
    };

    if (!activeLayout) {
        return <div>Loading dashboard...</div>;
    }

    return (
        <PageLayout onProjectSelect={handleProjectSelect}>
            {projectId && (
                <SuiteMenu
                    projectId={projectId as number}
                    onSuiteSelect={handleSuiteSelect}
                />
            )}
            <div className={`home-page ${projectId ? 'dashboard-with-sidebar' : ''}`}>
                <div className="dashboard-header">
                    <h2>{activeLayout.name}{projectId && ` - Project ${projectId}`}</h2>
                    <button className="btn btn-primary" onClick={() => setIsEditing(!isEditing)}>

                        {isEditing ? 'Done' : 'Edit Dashboard'}
                    </button>
                </div>

                {isEditing && (
                    <DashboardEditor onAddComponent={handleAddComponent} />
                )}

                <DashboardContainer
                    layout={activeLayout}
                    isEditing={isEditing}
                    onLayoutChange={updateGridLayout}
                    onRemoveComponent={removeComponent}
                />
            </div>
        </PageLayout>
    );
};

export default DashboardPage;
