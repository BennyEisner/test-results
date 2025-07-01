interface TableProps {
  bordered?: boolean;
  hover?: boolean;
  children: React.ReactNode;
  className?: string;
}

export const BootstrapTable = ({ 
  bordered = false, 
  hover = false, 
  children, 
  className = '' 
}: TableProps) => {
  const classes = [
    'table',
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
