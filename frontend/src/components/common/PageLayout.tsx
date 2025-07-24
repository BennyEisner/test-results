import React from 'react';
import BreadcrumbNavbar from './BreadcrumbNavbar';

interface PageLayoutProps {
  children: React.ReactNode;
  onProjectSelect?: (projectId: number) => void;
}

const PageLayout = ({ children, onProjectSelect }: PageLayoutProps) => {
  return (
    <div>
      <BreadcrumbNavbar onProjectSelect={onProjectSelect} />
      <main className="page-container">
        {children}
      </main>
    </div>
  );
};

export default PageLayout;
