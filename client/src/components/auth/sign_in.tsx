import { useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';

const SignIn: React.FC = () => {
  const [email, setEmail] = useState('');
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
    formData.append('email', email);
    formData.append('password', password);

    try {
      await axios.post('api/signin', formData).then(() => history.push('/dashboard'));
    } catch (err) {
      setMessage(err.response.data);
    }
  };

  return (
    <div>
      <h3>Sign In</h3>

      {message ? <p>{message}</p> : null}

      <form onSubmit={handleSubmit}>
        <label htmlFor='email'>email</label>
        <input id='email' type='email' onChange={ev => setEmail(ev.target.value)} />

        <label htmlFor='password'>password</label>
        <input id='password' type='password' onChange={ev => setPassword(ev.target.value)} />

        <button type='submit'>Sign in</button>
      </form>
    </div>
  );
};

export default SignIn;
