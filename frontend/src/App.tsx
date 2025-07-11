import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import ProjectsTable from './components/project/ProjectsTable';
import ProjectDetail from './components/project/ProjectDetail';
import SuiteDetail from './components/suite/SuiteDetail';
import BuildsTable from './components/build/BuildsTable';
import BuildDetail from './components/build/BuildDetail.tsx';
import HomePage from './components/page/HomePage';
import DashboardPage from './components/page/DashboardPage';
import PageLayout from './components/common/PageLayout';
import { AuthProvider } from './context/AuthContext';
import ProtectedRoute from './components/auth/ProtectedRoute';
import UserProfile from './components/auth/UserProfile';
import './styles/shared.css';
import './styles/tables.css';

function App() {
    return (
        <AuthProvider>
            <Router>
                <div className="app-container">
                    <Routes>
                        {/* Public routes */}
                        <Route path="/" element={<PageLayout><HomePage /></PageLayout>} />
                        
                        {/* Protected routes */}
                        <Route path="/dashboard" element={
                            <ProtectedRoute>
                                <DashboardPage />
                            </ProtectedRoute>
                        } />
                        <Route path="/projects" element={
                            <ProtectedRoute>
                                <PageLayout><ProjectsTable /></PageLayout>
                            </ProtectedRoute>
                        } />
                        <Route path="/projects/:projectId" element={
                            <ProtectedRoute>
                                <PageLayout><ProjectDetail /></PageLayout>
                            </ProtectedRoute>
                        } />
                        <Route path="/projects/:projectId/suites/:suiteId" element={
                            <ProtectedRoute>
                                <PageLayout><SuiteDetail /></PageLayout>
                            </ProtectedRoute>
                        } />
                        <Route
                            path="/projects/:projectId/suites/:suiteId/builds"
                            element={
                                <ProtectedRoute>
                                    <PageLayout><BuildsTableWrapper /></PageLayout>
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/projects/:projectId/suites/:suiteId/builds/:buildId"
                            element={
                                <ProtectedRoute>
                                    <PageLayout><BuildDetail /></PageLayout>
                                </ProtectedRoute>
                            }
                        />
                        <Route path="/profile" element={
                            <ProtectedRoute>
                                <PageLayout><UserProfile /></PageLayout>
                            </ProtectedRoute>
                        } />
                    </Routes>
                </div>
            </Router>
        </AuthProvider>
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
