import { createContext, useContext, useState, useEffect, ReactNode, useMemo } from 'react';
import { User } from '../types/auth';
import { authApi } from '../services/authApi';

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    error: string | null;
    login: (provider: string) => void;
    logout: () => void;
    refreshUser: () => Promise<void>;
    clearError: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
    children: ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const checkAuthStatus = async () => {
        try {
            setError(null);
            const userData = await authApi.getCurrentUser();
            setUser(userData);
        } catch (error) {
            // User not authenticated - this is expected for unauthenticated users
            setUser(null);
        } finally {
            setIsLoading(false);
        }
    };

    useEffect(() => {
        checkAuthStatus();
    }, []);

    const login = (provider: string) => {
        try {
            setError(null);
            authApi.beginOAuth2Auth(provider);
        } catch (err) {
            setError('Login failed. Please try again.');
        }
    };

    const logout = async () => {
        try {
            await authApi.logout();
            setUser(null);
            // The redirection is now handled in App.tsx
        } catch (error) {
            console.error('Logout error:', error);
        }
    };

    const refreshUser = async () => {
        await checkAuthStatus();
    };

    const clearError = () => {
        setError(null);
    };

    const value = useMemo(() => ({
        user,
        isAuthenticated: !!user,
        isLoading,
        error,
        login,
        logout,
        refreshUser,
        clearError
    }), [user, isLoading, error]);

    return (
        <AuthContext.Provider value={value}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = () => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};
