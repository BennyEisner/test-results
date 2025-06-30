import React, { useState, useEffect } from 'react';
import { fetchMostFailedTests } from '../services/api';
import {
    createColumnHelper,
    flexRender,
    getCoreRowModel,
    useReactTable,
} from '@tanstack/react-table';
import { Alert } from 'react-bootstrap';
import { MostFailedTest } from '../types';

interface MostFailedTestsTableProps {
    projectId: number;
}

const columnHelper = createColumnHelper<MostFailedTest>();

const columns = [
    columnHelper.accessor('name', {
        cell: (info) => info.getValue(),
        header: () => <span>Test Name</span>,
    }),
    columnHelper.accessor('classname', {
        cell: (info) => info.getValue(),
        header: () => <span>Class Name</span>,
    }),
    columnHelper.accessor('failure_count', {
        cell: (info) => info.getValue(),
        header: () => <span>Failure Count</span>,
    }),
];

const MostFailedTestsTable: React.FC<MostFailedTestsTableProps> = ({ projectId }) => {
    const [tests, setTests] = useState<MostFailedTest[]>([]);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const getTests = async () => {
            try {
                const data = await fetchMostFailedTests(projectId);
                setTests(data);
            } catch (err) {
                setError('Failed to fetch most failed tests');
            }
        };

        getTests();
    }, [projectId]);

    const table = useReactTable({
        data: tests,
        columns,
        getCoreRowModel: getCoreRowModel(),
    });

    if (error) {
        return <Alert variant="danger">{error}</Alert>;
    }

    return (
        <table className="table table-bordered table-hover">
            <thead>
                {table.getHeaderGroups().map(headerGroup => (
                    <tr key={headerGroup.id}>
                        {headerGroup.headers.map(header => (
                            <th key={header.id}>
                                {header.isPlaceholder
                                    ? null
                                    : flexRender(
                                        header.column.columnDef.header,
                                        header.getContext()
                                    )}
                            </th>
                        ))}
                    </tr>
                ))}
            </thead>
            <tbody>
                {table.getRowModel().rows.map(row => (
                    <tr key={row.id}>
                        {row.getVisibleCells().map(cell => (
                            <td key={cell.id}>
                                {flexRender(cell.column.columnDef.cell, cell.getContext())}
                            </td>
                        ))}
                    </tr>
                ))}
            </tbody>
        </table>
    );
};

export default MostFailedTestsTable;
