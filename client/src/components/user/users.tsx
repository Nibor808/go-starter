import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { IUser } from './interfaces';
import { useNavigate } from 'react-router-dom';

const Users: React.FC = () => {
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [users, setUsers] = useState<[IUser]>([
        {
            id: '',
            email: '',
            isActive: false,
            isAdmin: false,
        },
    ]);

    useEffect(() => {
        const id = setTimeout(() => {
            setError('');
        }, 2000);

        return () => clearTimeout(id);
    }, [error]);

    useEffect(() => {
        (async () => {
            try {
                const response = await axios.get('api/users');

                setUsers(response.data);
            } catch (err: any) {
                if (err.response.status === 401) {
                    setError(err.response.data);

                    setTimeout(() => {
                        return navigate('/');
                    }, 2000);
                }

                setError(err.response.data);
            }
        })();
    }, [navigate]);

    return (
        <>
            <p>All Users</p>
            {error ? <p className="error">ERROR: {error}</p> : null}

            {users.map((user) => (
                <ul key={user.id}>
                    <li>id: {user.id}</li>
                    <li>email: {user.email}</li>
                    <li>isActive: {String(user.isActive)}</li>
                    <li>isAdmin: {String(user.isAdmin)}</li>
                </ul>
            ))}
        </>
    );
};

export default Users;
