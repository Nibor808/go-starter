import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import axios from 'axios';

const ConfirmEmail: React.FC = () => {
    const [error, setError] = useState('');
    const [success, setSuccess] = useState('');
    const { token, userID } = useParams();
    const navigate = useNavigate();

    useEffect(() => {
        (async () => {
            try {
                const response = await axios.get(`http://localhost:3000/api/confirmemaildata/${token}/${userID}`);
                setSuccess(response.data);
                navigate('/signuppassword');
            } catch (err: any) {
                setError(err.response.data);
            }
        })();
    }, []);

    return (
        <div>
            <h3>Confirm Email</h3>

            <p>{error}</p>
            <p>{success}</p>
        </div>
    );
};

export default ConfirmEmail;
