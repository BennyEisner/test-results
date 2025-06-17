import { useParams, useNavigate } from 'react-router-dom';
//import BuildsTable from './BuildsTable';
import SuitesTable from './SuitesTable.tsx';

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
          Back to Projects
        </button>
        <h1>Project {projectId}</h1>
      </div>

      <SuitesTable projectId={projectId} />
    </div>
  );
};

export default ProjectDetail;
