//Boilerplate to be included on all pages 
import React from 'react';
import BreadcrumbNavbar from './BreadcrumbNavbar';

interface PageLayoutProps {
  children: React.ReactNode;
}

const PageLayout = ({ children }: PageLayoutProps) => {
  return (
    <div>
      <BreadcrumbNavbar />
      <main className="page-container">
        {children}
      </main>
    </div>
  );
};

export default PageLayout;
