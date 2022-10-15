import React, { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';

const SignUpPassword: React.FC = () => {
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');
    const history = useHistory();

    useEffect(() => {
        setTimeout(() => {
            setMessage('');
        }, 2000);
    }, [message]);

    const handleSubmit: React.FormEventHandler = async (ev: React.FormEvent) => {
        ev.preventDefault();

        const formData = new FormData();
        formData.append('password', password);

        try {
            await axios.post('api/signuppassword', formData);

            history.push('dashboard');
        } catch (err) {
            setMessage(err.response.data);
        }
    };

    return (
        <div>
            <h3>Sign Up</h3>

            <p>Now add a password.</p>
            {message ? <p>{message}</p> : null}

            <form onSubmit={handleSubmit} method="POST" className="v-form">
                <label htmlFor="password">password</label>
                <input
                    id="password"
                    type="password"
                    onChange={(ev: React.ChangeEvent<HTMLInputElement>) => setPassword(ev.target.value)}
                />

                <button type="submit">Continue</button>
            </form>
        </div>
    );
};

export default SignUpPassword;
