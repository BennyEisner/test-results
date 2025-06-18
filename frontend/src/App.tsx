import { BrowserRouter as Router, Routes, Route, useParams } from 'react-router-dom';
import ProjectsTable from './components/ProjectsTable';
import ProjectDetail from './components/ProjectDetail';
import SuiteDetail from './components/SuiteDetail';
import BuildsTable from './components/BuildsTable';
// import ExecutionsTable from './components/ExecutionsTable.tsx'; // Remove .tsx from import
import ExecutionsTable from './components/ExecutionsTable';


function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<ProjectsTable />} />
          <Route path="/projects" element={<ProjectsTable />} />
          <Route path="/projects/:projectId" element={<ProjectDetail />} />
          <Route path="/projects/:projectId/suites/:suiteId" element={<SuiteDetail />} />
          <Route 
            path="/projects/:projectId/suites/:suiteId/builds" 
            element={<BuildsTableWrapper />} 
          />
          <Route path="/builds/:buildId/executions" element={<ExecutionsTableWrapper />} />
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

// Wrapper component to extract params and pass them to ExecutionsTable
const ExecutionsTableWrapper = () => {
  const { buildId } = useParams<{ buildId: string }>();

  if (!buildId) {
    // Handle the case where params are not available
    return <div>Error: Missing build ID.</div>;
  }

  return <ExecutionsTable buildId={buildId} />;
};

export default App;
