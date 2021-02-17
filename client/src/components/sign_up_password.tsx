import { useState } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';

const SignUpPassword: React.FC = () => {
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const history = useHistory();

  const handleSubmit: React.FormEventHandler = async (ev: React.FormEvent) => {
    ev.preventDefault();

    const formData = new FormData();
    formData.append('password', password);

    try {
      const response = await axios.post('api/signuppassword', formData);

      setSuccess(response.data);
      history.push('dashboard');
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <h3>Sign Up</h3>

      <p>Now add a password.</p>
      <p>{error}</p>
      <p>{success}</p>

      <form onSubmit={handleSubmit} method='POST'>
        <label htmlFor='password'>password</label>
        <input id='password' type='password' onChange={ev => setPassword(ev.target.value)} />

        <button type='submit'>Continue</button>
      </form>
    </div>
  );
};

export default SignUpPassword;
