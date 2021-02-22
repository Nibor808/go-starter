import { useState, useEffect } from 'react';
import axios from 'axios';
import { IUser } from './interfaces';
import { useHistory } from 'react-router-dom';

const User: React.FC = () => {
  const history = useHistory();
  const [error, setError] = useState('');
  const [user, setUser] = useState<IUser>({
    id: '',
    email: '',
    isActive: false,
    isAdmin: false,
  });

  useEffect(() => {
    const timeout = setTimeout(() => {
      setError('');
    }, 2000);

    return () => clearTimeout(timeout);
  }, [error]);

  const handleUser = async () => {
    try {
      const response = await axios.get('api/user');

      setUser(response.data);
    } catch (err) {
      if (err.response.status === 401) {
        setError(err.response.data);

        setTimeout(() => {
          return history.push('/');
        }, 1500);
      }

      setError(err.response.data);
    }
  };

  return (
    <>
      <p>User From Cookie</p>
      {error ? <p>ERROR: {error}</p> : null}
      <ul>
        <li>{user.email}</li>
      </ul>

      <button onClick={handleUser}>User From Cookie</button>
    </>
  );
};

export default User;
