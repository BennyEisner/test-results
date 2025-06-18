import { useParams, useNavigate } from 'react-router-dom';
import BuildsTable from './BuildsTable';
//import SuitesTable from './SuitesTable.tsx';

const SuiteDetail = () => {
  const { suiteId, projectId} = useParams<{ suiteId: string, projectId: string}>();
  const navigate = useNavigate();

 if (!suiteId) {
    return <div className="error">Suite ID is required</div>;
  }

if (!projectId) {
    return <div className="error">Project ID is required</div>;
  }


  return (
    <div className="project-detail">
      <div className="project-header">
        <button 
          onClick={() => navigate(`/projects/${projectId}`)} 
          className="back-button"
        >
          Back to Suites
        </button>
        <h1>Suite {suiteId}</h1>
      </div>

      <BuildsTable  projectId={projectId} suiteId={suiteId}/>
    </div>
  );
};

export default SuiteDetail;
