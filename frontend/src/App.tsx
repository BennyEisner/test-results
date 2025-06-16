import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProjectsTable from './components/ProjectsTable';

function App() {
  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/" element={<ProjectsTable />} />
          {/* Add other routes here */}
        </Routes>
      </div>
    </Router>
  );
}

export default App;
