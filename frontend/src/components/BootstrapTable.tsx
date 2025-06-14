import React from 'react';

interface TableProps {
  striped?: boolean;
  bordered?: boolean;
  hover?: boolean;
  children: React.ReactNode;
  className?: string;
}

export const BootstrapTable: React.FC<TableProps> = ({ 
  striped = false, 
  bordered = false, 
  hover = false, 
  children, 
  className = '' 
}) => {
  const classes = [
    'table',
    striped ? 'table-striped' : '',
    bordered ? 'table-bordered' : '',
    hover ? 'table-hover' : '',
    className
  ].filter(Boolean).join(' ');

  return (
    <table className={classes}>
      {children}
    </table>
  );
};
