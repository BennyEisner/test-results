import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, useParams, useNavigate, useLocation } from 'react-router-dom';
import ProjectsTable from './components/project/ProjectsTable';
import ProjectDetail from './components/project/ProjectDetail';
import SuiteDetail from './components/suite/SuiteDetail';
import BuildsTable from './components/build/BuildsTable';
import BuildDetail from './components/build/BuildDetail.tsx';
import HomePage from './components/page/HomePage';
import DashboardPage from './components/page/DashboardPage';
import PageLayout from './components/common/PageLayout';
import { AuthProvider, useAuth } from './context/AuthContext';
import ProtectedRoute from './components/auth/ProtectedRoute';
import LoginPage from './components/auth/LoginPage';
import UserProfile from './components/auth/UserProfile';
import './styles/shared.css';
import './styles/tables.css';

const AppRoutes = () => {
    const { isAuthenticated, isLoading } = useAuth();
    const navigate = useNavigate();
    const location = useLocation();

    useEffect(() => {
        if (!isLoading) {
            if (isAuthenticated && location.pathname === '/login') {
                navigate('/dashboard');
            } else if (!isAuthenticated && location.pathname !== '/login') {
                // Optional: redirect to login if not authenticated and not on a public route
                // This might be too aggressive depending on desired UX
            }
        }
    }, [isAuthenticated, isLoading, navigate, location.pathname]);

    return (
        <Routes>
            {/* Public routes */}
            <Route path="/" element={<PageLayout><HomePage /></PageLayout>} />
            <Route path="/login" element={<PageLayout><LoginPage /></PageLayout>} />
            
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
    );
}

function App() {
    return (
        <AuthProvider>
            <Router>
                <div className="app-container">
                    <AppRoutes />
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
