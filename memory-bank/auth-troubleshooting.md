# Auth Troubleshooting Log

This document tracks the steps taken to diagnose and resolve the GitHub OAuth2 authentication issue.

## Problem: "you must select a provider"

The core issue is a `400 Bad Request` with the error "you must select a provider" when attempting to log in via GitHub. This indicates the `gothic` library cannot determine which provider to use from the HTTP request.

### Attempt 1: Fix Conflicting Initializations

- **Hypothesis**: A duplicate, incorrect Goth initialization in `container.go` was overwriting the correct session store configuration from `main.go`.
- **Action**: Removed the conflicting `FilesystemStore` initialization from `container.go`, ensuring the `CookieStore` from `main.go` is the single source of truth.
- **Result**: **Failed**. The error persisted. While cleaning up the duplicate initialization was correct, it was not the root cause of this specific error.

### Attempt 2: Override Provider Name Extraction

- **Hypothesis**: The default mechanism `gothic` uses to parse the provider name from the URL path is failing.
- **Action**: Overrode the default `gothic.GetProviderName` function in `main.go` with a more robust implementation that uses `r.PathValue("provider")` to reliably extract the provider from the URL.
- **Result**: **Failed**. The error persisted. This indicates the issue is not with the provider name extraction itself.

### Attempt 3: Manual URL Parsing and Diagnostic Logging

- **Hypothesis**: The issue may be related to how the router is handling the request, or some other subtle issue.
- **Action**:
    1.  Replaced the `r.PathValue`-based provider extraction with a manual URL parsing method in `main.go`. This is more resilient and not dependent on the router.
    2.  Added diagnostic logging to the `OAuth2Callback` handler to inspect the incoming request URL and headers.
- **Result**: **Partially Successful**. This resolved the "you must select a provider" error, but revealed a new issue.

### Attempt 4: Correct Environment Variable Mismatch

- **Hypothesis**: The application is not loading the GitHub client ID because of a mismatch between the environment variable name in the code (`GITHUB_KEY`) and the standard name (`GITHUB_CLIENT_ID`).
- **Action**:
    1.  Updated `api/cmd/server/main.go` to use `GITHUB_CLIENT_ID` when loading the GitHub credentials.
    2.  Updated `.env.example` to use `GITHUB_CLIENT_ID` for consistency.
- **Result**: **Successful**. This resolved the issue of the missing client ID.

## Final Diagnosis

The authentication flow is now working correctly from a code perspective. The final remaining error is `incorrect_client_credentials`. This indicates that the values for `GITHUB_CLIENT_ID` and `GITHUB_SECRET` in the `.env` file are incorrect. The application logic is sound, but the environment is misconfigured.
