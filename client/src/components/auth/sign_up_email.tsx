import { useState, useEffect } from 'react';
import axios from 'axios';

const SignUpEmail: React.FC = () => {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');

  useEffect(() => {
    setTimeout(() => {
      setMessage('');
    }, 2000);
  }, [message]);

  const handleSubmit: React.FormEventHandler = async (ev: React.FormEvent) => {
    ev.preventDefault();

    const formData = new FormData();
    formData.append('email', email);

    try {
      const response = await axios.post('api/signupemail', formData);

      setMessage(response.data);
    } catch (err) {
      setMessage(err.response.data);
    }
  };

  return (
    <div>
      <h3>Sign Up</h3>

      <p>Provide an email address.</p>
      <p>An email will be sent for you to confirm.</p>
      <p>Click on the link in the email to proceed.</p>

      {message ? <p>{message}</p> : null}

      <form onSubmit={handleSubmit} method='POST' className='v-form'>
        <label htmlFor='email'>email</label>
        <input id='email' type='email' value={email} onChange={ev => setEmail(ev.target.value)} />

        <button type='submit'>Send</button>
      </form>
    </div>
  );
};

export default SignUpEmail;
