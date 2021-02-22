import { IUser } from './interfaces';

interface IUserProps {
  user: IUser;
  onClick: () => any;
}

const User: React.FC<IUserProps> = ({ user, onClick }) => {
  return (
    <>
      <p>User From Cookie</p>
      <ul>
        <li>{user.email}</li>
      </ul>

      <button onClick={onClick}>User From Cookie</button>
    </>
  );
};

export default User;
