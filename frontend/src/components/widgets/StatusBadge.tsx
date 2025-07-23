import React from 'react';
import './StatusBadge.css';

interface StatusBadgeProps {
  status: 'success' | 'warning' | 'danger' | 'info' | 'neutral';
  children: React.ReactNode;
}

const StatusBadge: React.FC<StatusBadgeProps> = ({ status, children }) => {
  return (
    <span className={`status-badge status-badge-${status}`}>
      {children}
    </span>
  );
};

export default StatusBadge;
