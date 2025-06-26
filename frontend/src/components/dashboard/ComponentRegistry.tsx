import React from 'react';
import { ComponentType, ComponentProps } from '../../types/dashboard';
import BuildsTable from '../BuildsTable';
import ExecutionsSummary from '../ExecutionsSummary';
import BuildDoughnutChart from '../BuildDoughnutChart';
import { useExecutionsSummary } from '../../hooks/useExecutionsSummary';


interface ComponentRegistryProps {
    type: ComponentType;
    props: ComponentProps;
}

const ComponentRegistry = ({ type, props }: ComponentRegistryProps) => {

    const baseProps = { ...props, className: `dashboard-component dashboard-component--${type}`, };
    switch (type) {
        case 'builds-table':
            return (
                <div className="component-wrapper">
                    {props.title && <h3 className="component-title">{props.title}</h3>}
                    <BuildsTable {...baseProps} />
                </div>
            );
        case 'executions-summary':
            return (
                <div className="component-wrapper">
                    {props.title && <h3 className="component-title">{props.title}</h3>}
                    <ExecutionsSummary {...baseProps} />
                </div>
            );
        case 'build-chart':
            return (
                <div className="component-wrapper">
                    {props.title && <h3 className="component-title">{props.title}</h3>}
                    <BuildDoughnutChart {...baseProps} />
                </div>
            );

    }
