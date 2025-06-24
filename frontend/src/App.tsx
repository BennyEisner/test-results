import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import ProjectsTable from './components/ProjectsTable';
import ProjectDetail from './components/ProjectDetail';
import SuiteDetail from './components/SuiteDetail';
import BuildsTable from './components/BuildsTable';
import BuildDetail from './components/BuildDetail.tsx';
import HomePage from './components/HomePage';
import './styles/shared.css';

function App() {
    return (
        <Router>
            <div className="app-container">
                <Routes>
                    <Route path="/" element={<HomePage />} />
                    <Route path="/projects" element={<ProjectsTable />} />
                    <Route path="/projects/:projectId" element={<ProjectDetail />} />
                    <Route path="/projects/:projectId/suites/:suiteId" element={<SuiteDetail />} />
                    <Route
                        path="/projects/:projectId/suites/:suiteId/builds"
                        element={<BuildsTableWrapper />}
                    />
                    <Route
                        path="/projects/:projectId/suites/:suiteId/builds/:buildId"
                        element={<BuildDetail />}
                    />
                </Routes>
            </div>
        </Router>
    );
}
// Wrapper component to extract params and pass them to BuildsTable
const BuildsTableWrapper = () => {
    const { projectId, suiteId } = useParams<{ projectId: string; suiteId: string }>();

    if (!projectId || !suiteId) {
        // Handle the case where params are not available, though this shouldn't happen with a matched route
        return <div>Error: Missing project or suite ID.</div>;
    }

    return <BuildsTable projectId={projectId} suiteId={suiteId} />;
};



export default App;
