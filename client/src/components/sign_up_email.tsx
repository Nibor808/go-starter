import { useState } from 'react';
import axios from 'axios';

const SignUpEmail: React.FC = () => {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit: React.FormEventHandler = async (ev: React.FormEvent) => {
    ev.preventDefault();

    const formData = new FormData();
    formData.append('email', email);

    try {
      const response = await axios.post('api/signupemail', formData);

      setSuccess(response.data);
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <h3>Sign Up</h3>

      <p>Provide an email address.</p>
      <p>An email will be sent for you to confirm.</p>
      <p>Click on the link in the email to proceed.</p>

      <p>{error}</p>
      <p>{success}</p>

      <form onSubmit={handleSubmit} method='POST'>
        <label htmlFor='email'>email</label>
        <input id='email' type='email' value={email} onChange={ev => setEmail(ev.target.value)} />

        <button type='submit'>Send</button>
      </form>
    </div>
  );
};

export default SignUpEmail;
