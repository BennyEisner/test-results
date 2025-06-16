import { useParams, useNavigate } from 'react-router-dom';
import BuildsTable from './BuildsTable';

const ProjectDetail = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const navigate = useNavigate();

  if (!projectId) {
    return <div className="error">Project ID is required</div>;
  }

  return (
    <div className="project-detail">
      <div className="project-header">
        <button 
          onClick={() => navigate('/projects')} 
          className="back-button"
        >
          ← Back to Projects
        </button>
        <h1>Project {projectId}</h1>
      </div>

      <BuildsTable projectId={projectId} />
    </div>
  );
};

export default ProjectDetail;
