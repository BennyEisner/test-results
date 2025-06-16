import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import type { Project } from '../types';
import { fetchProjects } from '../services/api';
import './ProjectsTable.css';

const ProjectsTable = () => {
  const navigate = useNavigate();
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadProjects = async () => {
      try {
        setLoading(true);
        const data = await fetchProjects();
        setProjects(data);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load projects');
      } finally {
        setLoading(false);
      }
    };
    loadProjects();
  }, []);

  const handleProjectClick = (projectId: string | number) => {
    navigate(`/projects/${projectId}`);
  };

  if (loading) {
    return <div className="loading">Loading projects...</div>;
  }

  if (error) {
    return <div className="error">Error: {error}</div>;
  }

  return (
    <div>
      <table className="table table-striped table-bordered table-hover">
        <thead>
          <tr>
            <th>ID</th>
            <th>Project Name</th>
          </tr>
        </thead>
        <tbody>
          {projects.map((project) => (
            <tr 
              key={project.id}
              onClick={() => handleProjectClick(project.id)}
              style={{ cursor: 'pointer' }}
              className="clickable-row"
            >
              <td>{project.id}</td>
              <td>{project.name}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {projects.length === 0 && (
        <p className="no-data">No projects found.</p>
      )}
    </div>
  );
};

export default ProjectsTable;
