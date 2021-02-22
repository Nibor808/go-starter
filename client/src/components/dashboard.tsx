import Users from './user/users';
import User from './user/user';
import MyWebSocket from './websocket';

const Dashboard: React.FC = () => {
  return (
    <div>
      <h3>Dashboard</h3>

      <MyWebSocket />

      <Users />

      <User />
    </div>
  );
};

export default Dashboard;
