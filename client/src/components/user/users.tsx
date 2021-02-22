import { useState, useEffect } from 'react';
import axios from 'axios';
import { IUser } from './interfaces';
import { useHistory } from 'react-router-dom';

const Users: React.FC = () => {
  const history = useHistory();
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
    const timeout = setTimeout(() => {
      setError('');
    }, 2000);

    return () => clearTimeout(timeout);
  }, [error]);

  const handleUsers = async () => {
    try {
      const response = await axios.get('api/users');

      setUsers(response.data);
    } catch (err) {
      if (err.response.status === 401) return history.push('/');
      setError(err.response.data);
    }
  };

  return (
    <>
      <p>All Users</p>
      {error ? <p>ERROR: {error}</p> : null}

      <ul>
        {users.map(user => (
          <li key={user.id}>{user.email}</li>
        ))}
      </ul>

      <button onClick={handleUsers}>All Users</button>
    </>
  );
};

export default Users;
