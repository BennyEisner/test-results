import { useEffect } from 'react';
import { Responsive, WidthProvider } from 'react-grid-layout';
import { useDashboard } from '../../context/DashboardContext';
import { DashboardLayout, GridLayoutItem } from '../../types/dashboard';
import MemoizedComponentRegistry from './ComponentRegistry';
import BuildSelect from '../build/BuildSelect';
import 'react-grid-layout/css/styles.css';
import 'react-resizable/css/styles.css';
import '../../styles/dashboard.css';

const ResponsiveGridLayout = WidthProvider(Responsive);

interface DashboardContainerProps {
    layout: DashboardLayout;
    isEditing?: boolean;
    onLayoutChange?: (gridLayout: GridLayoutItem[]) => void;
    onRemoveComponent?: (componentId: string) => void;
}

const DashboardContainer = ({
    layout,
    isEditing = false,
    onLayoutChange,
    onRemoveComponent,
}: DashboardContainerProps) => {
    const { projectId, suiteId, buildId, setBuildId } = useDashboard();

    useEffect(() => {
        // Reset selected build when project or suite changes
        setBuildId(null);
    }, [projectId, suiteId, setBuildId]);

    const handleGridLayoutChange = (newGridLayout: any[]) => {
        if (onLayoutChange) {
            onLayoutChange(newGridLayout);
        }
    };

    const handleBuildSelect = (buildId: string | number) => {
        setBuildId(buildId);
    };

    const visibleComponents = layout.components.filter(comp => comp.visible);
    const visibleGridLayout = layout.gridLayout.filter(item => visibleComponents.some(c => c.id === item.i));

    return (
        <div>
            {projectId && suiteId && (
                <BuildSelect
                    projectId={projectId}
                    suiteId={suiteId}
                    onBuildSelect={handleBuildSelect}
                />
            )}
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
                                <div className="remove-btn-wrapper" onClick={() => onRemoveComponent && onRemoveComponent(component.id)}>
                                    <button
                                        className="remove-btn"
                                    >
                                        x
                                    </button>
                                </div>
                            </div>
                        )}
                        <div className="component-content">
                            <MemoizedComponentRegistry
                                type={component.type}
                                props={component.props}
                                projectId={projectId ?? undefined}
                                suiteId={suiteId ?? undefined}
                                buildId={buildId}
                            />
                        </div>
                    </div>
                ))}
            </ResponsiveGridLayout>
        </div>
    );
};

export default DashboardContainer;
