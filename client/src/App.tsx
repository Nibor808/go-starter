import React, { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

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
        try {
            const response = await axios.get('api/signout');
            setError(response.data);
            navigate('/');
        } catch (err: any) {
            if (err.response.status === 401) {
                setError(err.response.data);

                setTimeout(() => {
                    return navigate('/');
                }, 1500);
            }

            setError(err.response.data);
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
