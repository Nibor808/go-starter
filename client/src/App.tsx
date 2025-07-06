import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { requestGet } from './components/request/get';

export const App: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [error, setError] = useState('');
    const navigate = useNavigate();

    useEffect(() => {
        const id = setTimeout(() => {
            setError('');
        }, 3000);

        return () => clearTimeout(id);
    }, [error]);

    const handleSignOut = async () => {
        const didSignOut = await requestGet({ url: 'api/signout' });

        if (didSignOut) {
            navigate('/');
        } else {
            setError('Could not sign out');
        }
    };

    return (
        <div className="app">
            <header className="app-header">
                <h1>Welcome to Go Starter!</h1>

                <div className="app-menu">
                    <button onClick={handleSignOut}>Sign Out</button>
                </div>
            </header>

            {error ? <p className="error">{error}</p> : null}

            {children}
        </div>
    );
};
