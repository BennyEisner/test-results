import ProjectsTable from './components/ProjectsTable'
import './App.css'

function App() {
  return (
    <div className="app">
      <header className="app-header">
        <h1>Test Results Dashboard</h1>
      </header>
      <main>
        <ProjectsTable />
      </main>
    </div>
  )
}

export default App
