interface TableProps {
  striped?: boolean;
  bordered?: boolean;
  hover?: boolean;
  children: React.ReactNode;
  className?: string;
}

export const BootstrapTable = ({ 
  striped = false, 
  bordered = false, 
  hover = false, 
  children, 
  className = '' 
}: TableProps) => {
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
