import React, { useState, useEffect, useRef } from 'react';
import { search } from '../../services/api';
import type { SearchResult } from '../../types';
import './SearchBar.css';

interface SearchBarProps {
    onResultSelect?: (result: SearchResult) => void;
}

const SearchBar = ({ onResultSelect }: SearchBarProps) => {
    const [query, setQuery] = useState('');
    const [results, setResults] = useState<SearchResult[]>([]);
    const [isLoading, setIsLoading] = useState(false);
    const [showDropdown, setShowDropdown] = useState(false);
    const [selectedIndex, setSelectedIndex] = useState(-1);
    const searchRef = useRef<HTMLDivElement>(null);
    const inputRef = useRef<HTMLInputElement>(null);

    useEffect(() => {
        const delayedSearch = setTimeout(async () => {
            if (query.trim()) {
                setIsLoading(true);
                try {
                    const searchResults = await search(query);
                    setResults(Array.isArray(searchResults) ? searchResults : []);
                    setShowDropdown(true);
                } catch (error) {
                    console.error('Search failed:', error);
                    setResults([]);
                } finally {
                    setIsLoading(false);
                }
            } else {
                setResults([]);
                setShowDropdown(false);
            }
        }, 300);

        return () => clearTimeout(delayedSearch);
    }, [query]);

    useEffect(() => {
        const handleClickOutside = (event: MouseEvent) => {
            if (searchRef.current && !searchRef.current.contains(event.target as Node)) {
                setShowDropdown(false);
                setSelectedIndex(-1);
            }
        };

        document.addEventListener('mousedown', handleClickOutside);
        return () => document.removeEventListener('mousedown', handleClickOutside);
    }, []);

    const handleKeyDown = (e: React.KeyboardEvent) => {
        if (!showDropdown || results.length === 0) return;

        switch (e.key) {
            case 'ArrowDown':
                e.preventDefault();
                setSelectedIndex(prev => (prev < results.length - 1 ? prev + 1 : 0));
                break;
            case 'ArrowUp':
                e.preventDefault();
                setSelectedIndex(prev => (prev > 0 ? prev - 1 : results.length - 1));
                break;
            case 'Enter':
                e.preventDefault();
                if (selectedIndex >= 0 && selectedIndex < results.length) {
                    handleResultClick(results[selectedIndex]);
                }
                break;
            case 'Escape':
                setShowDropdown(false);
                setSelectedIndex(-1);
                inputRef.current?.blur();
                break;
        }
    };

    const handleResultClick = (result: SearchResult) => {
        setQuery(result.name);
        setShowDropdown(false);
        setSelectedIndex(-1);
        if (onResultSelect) {
            onResultSelect(result);
        }
    };


    const getTypeLabel = (type: string) => {
        switch (type) {
            case 'project':
                return 'Project';
            case 'test_suite':
                return 'Test Suite';
            case 'build':
                return 'Build';
            case 'test_case':
                return 'Test Case';
            default:
                return type;
        }
    };

    return (
        <div className="search-bar" ref={searchRef}>
            <div className="search-input-container">
                <input
                    ref={inputRef}
                    type="text"
                    className="search-input"
                    placeholder="Search projects, suites, builds, or test cases..."
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    onKeyDown={handleKeyDown}
                    onFocus={() => {
                        if (results.length > 0) {
                            setShowDropdown(true);
                        }
                    }}
                />
                {isLoading && <div className="search-loading">🔍</div>}
            </div>

            {showDropdown && (
                <div className="search-dropdown">
                    {results.length === 0 && !isLoading && query.trim() && (
                        <div className="search-no-results">No results found</div>
                    )}
                    {results.map((result, index) => (
                        <div
                            key={`${result.type}-${result.id}`}
                            className={`search-result-item ${index === selectedIndex ? 'selected' : ''}`}
                            onClick={() => handleResultClick(result)}
                            onMouseEnter={() => setSelectedIndex(index)}
                        >
                            <div className="search-result-content">
                                <div className="search-result-name">{result.name}</div>
                                <div className="search-result-type">{getTypeLabel(result.type)}</div>
                            </div>
                        </div>
                    ))}
                </div>
            )}
        </div>
    );
};

export default SearchBar;
