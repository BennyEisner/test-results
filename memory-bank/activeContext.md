# Active Context

## Current Work Focus - CRITICAL ISSUE
**SYSTEM IS BROKEN** - The dashboard chart functionality is completely non-functional due to recent backend changes. All chart API endpoints are returning 400 Bad Request errors.

## Recent Changes That Broke The System
**Breaking Changes Made:**
1. **Modified HTTP Handler Parameter Parsing:** Changed `api/internal/dashboard/infrastructure/http/http_handler.go` to read `project_id` from query parameters instead of path parameters, breaking the existing API contract.
2. **Updated Service Layer Signatures:** Modified `DashboardService.GetChartData()` interface to accept `suiteID` and `buildID` parameters.
3. **Updated Repository Layer:** Modified `BuildTestCaseExecutionRepository.GetChartData()` to accept additional filtering parameters.
4. **Frontend API Changes:** Updated `dashboardApi.getChartData()` to send `suiteId` and `buildId` as query parameters.
5. **Component Registry Updates:** Added configuration fields for static charts in `ComponentRegistry.tsx`.

## State Before Breaking Changes
- Dashboard was functional with basic chart rendering
- Charts displayed data correctly for project-level filtering
- API endpoints worked with path-based project ID: `/api/dashboard/projects/{projectID}/chart/{chartType}`
- No suite or build-level filtering was implemented

## Current Error State
**API Errors:**
- `GET /api/dashboard/projects/1/chart/build-duration` returns 400 Bad Request
- `GET /api/dashboard/projects/1/chart/build-duration?suite_id=3` returns 400 Bad Request  
- `GET /api/dashboard/projects/1/chart/build-duration?suite_id=3&build_id=3` returns 400 Bad Request

## Root Cause Analysis
**Primary Issue:** The HTTP handler now expects `project_id` as a query parameter but the frontend is still sending it as a path parameter. This creates a mismatch where:
- Frontend sends: `/api/dashboard/projects/1/chart/build-duration`
- Backend expects: `/api/dashboard/chart/build-duration?project_id=1`

**Secondary Issues:**
1. Route configuration may not match the new parameter expectations
2. Service layer signature changes propagated through entire backend stack
3. Frontend and backend API contracts are now misaligned

## Immediate Investigation Needed
1. **Check Route Configuration:** Verify how the dashboard routes are configured in the main server
2. **API Contract Alignment:** Determine if we should revert the HTTP handler changes or update the frontend API calls
3. **Parameter Parsing:** Ensure the HTTP handler correctly reads parameters from the expected locations
4. **Service Layer Compatibility:** Verify all service implementations match the updated interfaces

## Recovery Strategy Options
1. **Revert Backend Changes:** Roll back all backend modifications to restore functionality
2. **Fix Parameter Parsing:** Update HTTP handler to read project_id from path while maintaining new query parameters
3. **Update Frontend:** Modify frontend to match new backend API contract
4. **Hybrid Approach:** Maintain backward compatibility while adding new functionality

## Next Steps & Issues
- **CRITICAL:** Restore basic chart functionality before implementing new features
- **Investigate:** Root cause of 400 Bad Request errors
- **Decide:** Whether to revert changes or fix forward
- **Test:** Ensure basic dashboard functionality works before adding enhancements
- **Document:** Proper API contracts for future development
