import { IUser } from './interfaces';

interface IUsersProps {
  users: [IUser];
  onClick: () => any;
}

const Users: React.FC<IUsersProps> = ({ users, onClick }) => {
  return (
    <>
      <p>All Users</p>

      <ul>
        {users.map(user => (
          <li key={user.id}>{user.email}</li>
        ))}
      </ul>

      <button onClick={onClick}>All Users</button>
    </>
  );
};

export default Users;
