import { useState, useEffect } from 'react';
import axios from 'axios';

interface User {
  Id: string;
  Email: string;
  Active: boolean;
  Admin: boolean;
}

const Dashboard: React.FC = () => {
  const [users, setUsers] = useState<[User]>([
    {
      Id: '',
      Email: '',
      Active: false,
      Admin: false,
    },
  ]);
  const [error, setError] = useState('');

  useEffect(() => {
    handleUsers();
  }, []);

  const handleUsers = async () => {
    try {
      const response = await axios.get('api/users');

      console.log(response);
      setUsers(response.data);
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <h3>Dashboard</h3>
      <p>{error}</p>

      <ul>
        {users.map(user => (
          <li key={user.Id}>{user.Email}</li>
        ))}
      </ul>

      <button onClick={handleUsers}>All Users</button>
      <button>User From Cookie</button>
    </div>
  );
};

export default Dashboard;
