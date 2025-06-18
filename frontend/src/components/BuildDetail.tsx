import { useParams, useNavigate } from 'react-router-dom';
// import BuildsTable from './BuildsTable'; // Unused import
// import SuitesTable from './SuitesTable.tsx'; // Unused import
import ExecutionsTable from './ExecutionsTable'; // Added import

const BuildDetail = () => {
  const { buildId, suiteId, projectId} = useParams<{ buildId: string, suiteId: string, projectId: string}>();
  const navigate = useNavigate();

 if (!suiteId) {
    return <div className="error">Suite ID is required</div>;
  }

if (!projectId) {
    return <div className="error">Project ID is required</div>;
  }

 if (!buildId) {
    // Corrected error message
    return <div className="error">Build ID is required</div>; 
  }

  return (
    <div className="project-detail">
      <div className="project-header">
        <button 
          // Corrected navigation path
          onClick={() => navigate(`/builds/${buildId}`)} 
          className="back-button"
        >
          Back to Builds 
        </button>
        <h1>Build {buildId}</h1> {/* Corrected heading */}
      </div>

      <ExecutionsTable buildId={buildId}/>
    </div>
  );
};

export default BuildDetail; // Corrected export
