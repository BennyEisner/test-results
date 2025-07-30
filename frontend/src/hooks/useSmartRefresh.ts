import { useState, useEffect, useRef } from 'react';

type RefreshTrigger = 'project' | 'suite' | 'build';

interface SmartRefreshOptions<T> {
    projectId?: string | number;
    suiteId?: string | number;
    buildId?: string | number | null;
    isStatic?: boolean;
    staticProjectId?: string | number;
    staticSuiteId?: string | number;
    staticBuildId?: string | number;
    limit?: number;
    fetcher: (
        projectId: number,
        suiteId?: number,
        buildId?: number,
        limit?: number,
        signal?: AbortSignal
    ) => Promise<T>;
    refreshOn: RefreshTrigger[];
}

interface SmartRefreshResult<T> {
    data: T | null;
    error: Error | null;
    isLoading: boolean;
}

export const useSmartRefresh = <T>({
    projectId,
    suiteId,
    buildId,
    isStatic,
    staticProjectId,
    staticSuiteId,
    staticBuildId,
    fetcher,
    refreshOn,
    limit,
}: SmartRefreshOptions<T>): SmartRefreshResult<T> => {
    const [data, setData] = useState<T | null>(null);
    const [error, setError] = useState<Error | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const isInitialMount = useRef(true);

    const effectiveProjectId = isStatic ? staticProjectId : projectId;
    const effectiveSuiteId = isStatic ? staticSuiteId : suiteId;
    const effectiveBuildId = isStatic ? staticBuildId : buildId;

    const prevContext = useRef({ effectiveProjectId, effectiveSuiteId, effectiveBuildId });

    useEffect(() => {
        const abortController = new AbortController();
        const fetchId = Date.now();
        console.log(`[${fetchId}] useEffect triggered`, {
            effectiveProjectId,
            effectiveSuiteId,
            effectiveBuildId,
        });

        const fetchData = async () => {
            console.log(`[${fetchId}] fetchData called`, {
                effectiveProjectId,
                effectiveSuiteId,
                effectiveBuildId,
            });
            if (!effectiveProjectId) {
                console.log(`[${fetchId}] No project ID, skipping fetch.`);
                setIsLoading(false);
                return;
            }
            setIsLoading(true);
            console.log(`[${fetchId}] Fetching data...`);
            try {
                const result = await fetcher(
                    Number(effectiveProjectId),
                    effectiveSuiteId ? Number(effectiveSuiteId) : undefined,
                    effectiveBuildId ? Number(effectiveBuildId) : undefined,
                    limit,
                    abortController.signal
                );
                if (!abortController.signal.aborted) {
                    console.log(`[${fetchId}] Fetch successful`, result);
                    setData(result);
                }
            } catch (err) {
                if (!abortController.signal.aborted) {
                    console.error(`[${fetchId}] Fetch error`, err);
                    setError(err as Error);
                }
            } finally {
                if (!abortController.signal.aborted) {
                    console.log(`[${fetchId}] Fetch complete, setting loading to false.`);
                    setIsLoading(false);
                }
            }
        };

        const shouldRefresh = () => {
            if (isInitialMount.current) {
                console.log(`[${fetchId}] shouldRefresh: initial mount`);
                return true;
            }

            const prev = prevContext.current;
            console.log(`[${fetchId}] shouldRefresh: comparing contexts`, {
                prev,
                current: { effectiveProjectId, effectiveSuiteId, effectiveBuildId },
            });

            if (
                refreshOn.includes('project') &&
                prev.effectiveProjectId !== effectiveProjectId
            ) {
                console.log(`[${fetchId}] shouldRefresh: project changed`);
                return true;
            }
            if (
                refreshOn.includes('suite') &&
                prev.effectiveSuiteId !== effectiveSuiteId
            ) {
                console.log(`[${fetchId}] shouldRefresh: suite changed`);
                return true;
            }
            if (
                refreshOn.includes('build') &&
                prev.effectiveBuildId !== effectiveBuildId
            ) {
                console.log(`[${fetchId}] shouldRefresh: build changed`);
                return true;
            }
            console.log(`[${fetchId}] shouldRefresh: no change detected`);
            return false;
        };

        if (shouldRefresh()) {
            fetchData();
            isInitialMount.current = false;
        }

        prevContext.current = { effectiveProjectId, effectiveSuiteId, effectiveBuildId };

        return () => {
            console.log(`[${fetchId}] Cleanup: aborting fetch`);
            abortController.abort();
        };
    }, [effectiveProjectId, effectiveSuiteId, effectiveBuildId, fetcher, limit, refreshOn]);

    return { data, error, isLoading };
};
