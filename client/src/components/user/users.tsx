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
    const id = setTimeout(() => {
      setError('');
    }, 2000);

    return () => clearTimeout(id);
  }, [error]);

  useEffect(() => {
    const handleUsers = async () => {
      try {
        const response = await axios.get('api/users');

        setUsers(response.data);
      } catch (err) {
        if (err.response.status === 401) {
          setError(err.response.data);

          setTimeout(() => {
            return history.push('/');
          }, 2000);
        }

        setError(err.response.data);
      }
    };

    handleUsers();
  }, [history]);

  return (
    <>
      <p>All Users</p>
      {error ? <p className='error'>ERROR: {error}</p> : null}

      {users.map(user => (
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
