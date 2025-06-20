import { Table, Spinner, Alert, Badge } from 'react-bootstrap';
import type { TestCaseExecution } from '../types';

interface ExecutionsTableProps {
    executions: TestCaseExecution[];
    loading: boolean;
}

const ExecutionsTable = ({ executions, loading }: ExecutionsTableProps) => {
    if (loading) {
        return (
            <div className="d-flex justify-content-center align-items-center" style={{ height: '60vh' }}>
                <Spinner animation="border" role="status">
                    <span className="visually-hidden">Loading executions...</span>
                </Spinner>
            </div>
        );
    }

    const getStatusBadge = (status: string, hasFailure: boolean) => {
        const upperStatus = status?.toUpperCase();
        if (hasFailure || upperStatus === 'FAILED') {
            return <Badge bg="danger">{status}</Badge>;
        }
        if (upperStatus === 'PASSED') {
            return <Badge bg="success">{status}</Badge>;
        }
        if (upperStatus === 'SKIPPED') {
            return <Badge bg="warning" text="dark">{status}</Badge>;
        }
        return <Badge bg="secondary">{status}</Badge>;
    };

    return (
        <div className="py-3">
            <h2 className="mb-3">Test Executions</h2>
            <Table striped bordered hover responsive>
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
                                {getStatusBadge(execution.status, !!execution.failure)}
                                {execution.failure && (
                                    <span 
                                        title={`Failure: ${execution.failure.message || 'No message'}`}
                                        className="ms-2 text-danger"
                                        style={{ cursor: 'help' }}
                                    >
                                        ⚠️
                                    </span>
                                )}
                            </td>
                            <td className="font-monospace">{execution.execution_time}s</td>
                            <td className="text-muted">{new Date(execution.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </Table>
            {executions.length === 0 && !loading && (
                <Alert variant="info" className="mt-3">No executions found for this build.</Alert>
            )}
        </div>
    );
};

export default ExecutionsTable;
