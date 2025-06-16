import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProjectsTable from './components/ProjectsTable';
import ProjectDetail from './components/ProjectDetail';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<ProjectsTable />} />
          <Route path="/projects" element={<ProjectsTable />} />
          <Route path="/projects/:projectId" element={<ProjectDetail />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
