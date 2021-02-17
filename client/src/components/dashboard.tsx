import { useState } from 'react';
import axios from 'axios';

interface User {
  id: string;
  email: string;
  isActive: boolean;
  isAdmin: boolean;
}

const Dashboard: React.FC = () => {
  const [users, setUsers] = useState<[User]>([
    {
      id: '',
      email: '',
      isActive: false,
      isAdmin: false,
    },
  ]);
  const [user, setUser] = useState<User>({
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
      <p>{error}</p>

      <p>All Users</p>
      <ul>
        {users.map(user => (
          <li key={user.id}>{user.email}</li>
        ))}
      </ul>

      <p>User From Cookie</p>
      <ul>
        <li>{user.email}</li>
      </ul>

      <button onClick={handleUsers}>All Users</button>
      <button onClick={handleUser}>User From Cookie</button>
    </div>
  );
};

export default Dashboard;
