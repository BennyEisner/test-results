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
        limit?: number
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
        const fetchData = async () => {
            if (!effectiveProjectId) {
                setIsLoading(false);
                return;
            }
            setIsLoading(true);
            try {
                const result = await fetcher(
                    Number(effectiveProjectId),
                    effectiveSuiteId ? Number(effectiveSuiteId) : undefined,
                    effectiveBuildId ? Number(effectiveBuildId) : undefined,
                    limit
                );
                setData(result);
            } catch (err) {
                setError(err as Error);
            } finally {
                setIsLoading(false);
            }
        };

        const shouldRefresh = () => {
            if (isInitialMount.current) {
                return true;
            }

            const prev = prevContext.current;
            if (
                refreshOn.includes('project') &&
                prev.effectiveProjectId !== effectiveProjectId
            ) {
                return true;
            }
            if (
                refreshOn.includes('suite') &&
                prev.effectiveSuiteId !== effectiveSuiteId
            ) {
                return true;
            }
            if (
                refreshOn.includes('build') &&
                prev.effectiveBuildId !== effectiveBuildId
            ) {
                return true;
            }
            return false;
        };

        if (shouldRefresh()) {
            fetchData();
            isInitialMount.current = false;
        }

        prevContext.current = { effectiveProjectId, effectiveSuiteId, effectiveBuildId };
    }, [effectiveProjectId, effectiveSuiteId, effectiveBuildId, fetcher, limit, refreshOn]);

    return { data, error, isLoading };
};
