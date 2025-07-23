import React from 'react';
import './MetricCard.css';

interface MetricCardProps {
  title?: string;
  value: string | number;
  change?: string;
  changeType?: 'increase' | 'decrease' | 'neutral';
}

const MetricCard: React.FC<MetricCardProps> = ({ title, value, change, changeType = 'neutral' }) => {
  const getChangeColor = () => {
    if (changeType === 'increase') return 'metric-change-increase';
    if (changeType === 'decrease') return 'metric-change-decrease';
    return 'metric-change-neutral';
  };

  return (
    <div className="metric-card">
      <h3 className="metric-title">{title}</h3>
      <div className="metric-value">{value}</div>
      {change && (
        <div className={`metric-change ${getChangeColor()}`}>
          {changeType === 'increase' ? '▲' : '▼'} {change}
        </div>
      )}
    </div>
  );
};

export default MetricCard;
