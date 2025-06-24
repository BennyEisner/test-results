// ExecutionsSummary.tsx
import { useMemo } from 'react';
import { Card, Row, Col, Spinner } from 'react-bootstrap';
import type { TestCaseExecution } from '../types';

interface ExecutionsSummaryProps {
  executions: TestCaseExecution[];
  loading: boolean;
}

const ExecutionsSummary = ({ executions, loading }: ExecutionsSummaryProps) => {
  // Calculate execution statistics only when the executions data changes
  const stats = useMemo(() => {
    if (loading || executions.length === 0) {
      return {
        total: 0,
        passed: 0,
        failed: 0,
        skipped: 0,
        passRate: 0,
        avgTime: 0
      };
    }

    // Count total executions
    const total = executions.length;
    
    // Count executions by status
    const passed = executions.filter(e => e.status?.toUpperCase() === 'PASSED' && e.failure == null).length;
    const failed = executions.filter(e => 
      e.status?.toUpperCase() === 'FAILED' || e.failure != null
    ).length;
    const skipped = executions.filter(e => e.status?.toUpperCase() === 'SKIPPED').length;
    
    // Calculate pass rate percentage
    const passRate = total > 0 ? Math.round((passed / total) * 100) : 0;
    
    // Calculate average execution time (in seconds)
    const totalTime = executions.reduce((sum, e) => 
      sum + (e.execution_time || 0), 0
    );
    const avgTime = total > 0 ? (totalTime / total) : 0;
    
    return { total, passed, failed, skipped, passRate, avgTime };
  }, [executions, loading]);

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center my-3">
        <Spinner animation="border" role="status" variant="primary" />
        <span className="ms-2">Analyzing test results...</span>
      </div>
    );
  }

  if (stats.total === 0) {
    return (
        <div className="my-3 text-center text-muted">
            No execution data to summarize.
        </div>
    );
  }

  // Render the summary metrics
  return (
    <div>
      <Row xs={1} sm={2} md={3} lg={6} className="g-3">
        <Col>
          <Card text="white" bg="primary" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.total}</Card.Title>
              <Card.Text>Total Tests</Card.Text>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card text="white" bg="success" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.passed}</Card.Title>
              <Card.Text>Passed</Card.Text>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card text="white" bg="danger" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.failed}</Card.Title>
              <Card.Text>Failed</Card.Text>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card text="white" bg="secondary" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.skipped}</Card.Title>
              <Card.Text>Skipped</Card.Text>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card text="white" bg="primary" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.passRate}%</Card.Title>
              <Card.Text>Pass Rate</Card.Text>
            </Card.Body>
          </Card>
        </Col>
        <Col>
          <Card text="white" bg="primary" className="h-100">
            <Card.Body className="text-center">
              <Card.Title as="h4">{stats.avgTime.toFixed(2)}s</Card.Title>
              <Card.Text>Avg. Time</Card.Text>
            </Card.Body>
          </Card>
        </Col>
      </Row>
    </div>
  );
};

export default ExecutionsSummary;
