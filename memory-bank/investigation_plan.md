# Investigation Plan: Dashboard Chart Issues

## 1. Static Charts Stuck on Loading

### Objective
Identify why the static charts are not rendering and remain in a "loading" state.

### Steps
1.  **Debug `useSmartRefresh` Hook:**
    -   File: `frontend/src/hooks/useSmartRefresh.ts`
    -   Action: Add logging to trace the `isLoading` state and see if it is being updated correctly.
2.  **Debug `DataChart` Component:**
    -   File: `frontend/src/components/widgets/DataChart.tsx`
    -   Action: Add logging to see if it is re-rendering when the data is available.

## 2. Missing "Limit" Option for Dynamic Charts

### Objective
Identify why the "limit" option is not available for dynamic charts.

### Steps
1.  **Inspect `ComponentRegistry`:**
    -   File: `frontend/src/components/dashboard/ComponentRegistry.tsx`
    -   Action: Review the configuration for the dynamic chart components to ensure that the "limit" option is correctly defined.
2.  **Inspect `ComponentConfigModal`:**
    -   File: `frontend/src/components/dashboard/ComponentConfigModal.tsx`
    -   Action: Review the `ComponentConfigModal` to ensure that it is correctly displaying the "limit" option.
