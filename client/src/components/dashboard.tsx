import { useState, useEffect } from 'react';
import { IUser } from './user/interfaces';
import Users from './user/users';
import User from './user/user';
import MyWebSocket from './websocket';
import axios from 'axios';

const Dashboard: React.FC = () => {
  const [users, setUsers] = useState<[IUser]>([
    {
      id: '',
      email: '',
      isActive: false,
      isAdmin: false,
    },
  ]);

  const [user, setUser] = useState<IUser>({
    id: '',
    email: '',
    isActive: false,
    isAdmin: false,
  });
  const [error, setError] = useState('');

  const handleUsers = async () => {
    try {
      const response = await axios.get('api/users');

      setUsers(response.data);
    } catch (err) {
      setError(err.response.data);
    }
  };

  useEffect(() => {
    setTimeout(() => {
      setError('');
    }, 2000);
  }, [error]);

  const handleUser = async () => {
    try {
      const response = await axios.get('api/user');

      setUser(response.data);
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <h3>Dashboard</h3>
      {error ? <p>ERROR: {error}</p> : null}

      <MyWebSocket />

      <Users users={users} onClick={handleUsers} />

      <User user={user} onClick={handleUser} />
    </div>
  );
};

export default Dashboard;
