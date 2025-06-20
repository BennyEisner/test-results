import type { TestCaseExecution } from '../types';

interface ExecutionsTableProps {
    executions: TestCaseExecution[];
    loading: boolean;
}

const ExecutionsTable = ({ executions, loading }: ExecutionsTableProps) => {
    if (loading) {
        return <div className="loading">Loading executions...</div>;
    }

    return (
        <div>
            <h2>Test Executions</h2>
            <table className="table table-striped table-bordered table-hover">
                <thead>
                    <tr>
                        <th>Execution ID</th>
                        <th>Test Case Name</th>
                        <th>Status</th>
                        <th>Execution Time</th>
                        <th>Created At</th>
                    </tr>
                </thead>
                <tbody>
                    {executions.map((execution) => (
                        <tr key={execution.id} className={execution.failure ? 'table-danger' : ''}>
                            <td>#{execution.id}</td>
                            <td>{execution.test_case_name || `Test Case ${execution.test_case_id}`}</td>
                            <td>
                                {execution.status}
                                {execution.failure && (
                                    <span title={`Failure: ${execution.failure.message || 'No message'}`}
                                        style={{ marginLeft: '8px', color: 'red', cursor: 'help' }}>
                                        ⚠️
                                    </span>
                                )}
                            </td>
                            <td>{execution.execution_time}s</td>
                            <td>{new Date(execution.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
            {executions.length === 0 && (
                <p className="no-data">No executions found for this build.</p>
            )}
        </div>
    );
};
export default ExecutionsTable;
