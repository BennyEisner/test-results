import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import ProjectsTable from './components/ProjectsTable';
import ProjectDetail from './components/ProjectDetail';
import SuiteDetail from './components/SuiteDetail';
import BuildsTable from './components/BuildsTable';
import BuildDetail from './components/BuildDetail.tsx';
import HomePage from './components/HomePage';
import DashboardPage from './components/DashboardPage';
import PageLayout from './components/PageLayout';
import './styles/shared.css';
import './styles/tables.css';

function App() {
    return (
        <Router>
            <div className="app-container">
                <Routes>
                    <Route path="/" element={<PageLayout><HomePage /></PageLayout>} />
                    <Route path="/dashboard" element={<PageLayout><DashboardPage /></PageLayout>} />
                    <Route path="/projects" element={<PageLayout><ProjectsTable /></PageLayout>} />
                    <Route path="/projects/:projectId" element={<PageLayout><ProjectDetail /></PageLayout>} />
                    <Route path="/projects/:projectId/suites/:suiteId" element={<PageLayout><SuiteDetail /></PageLayout>} />
                    <Route
                        path="/projects/:projectId/suites/:suiteId/builds"
                        element={<PageLayout><BuildsTableWrapper /></PageLayout>}
                    />
                    <Route
                        path="/projects/:projectId/suites/:suiteId/builds/:buildId"
                        element={<PageLayout><BuildDetail /></PageLayout>}
                    />
                </Routes>
            </div>
        </Router>
    );
}

// Wrapper component to extract params and pass them to BuildsTable
const BuildsTableWrapper = () => {
    const { projectId, suiteId } = useParams<{ projectId: string; suiteId:string }>();

    if (!projectId || !suiteId) {
        // Handle the case where params are not available, though this shouldn't happen with a matched route
        return <div>Error: Missing project or suite ID.</div>;
    }

    return <BuildsTable projectId={projectId} suiteId={suiteId} />;
};



export default App;
