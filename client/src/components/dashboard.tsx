import { useState } from 'react';
import Users from './user/users';
import User from './user/user';
import MyWebSocket from './websocket';

const Dashboard: React.FC = () => {
  const [currentView, setCurrentView] = useState('');

  return (
    <>
      <h3>Dashboard</h3>
      <div className='app-sub-menu'>
        <button onClick={() => setCurrentView('currentUser')}>Current User</button>
        <button onClick={() => setCurrentView('allUsers')}>All Users</button>
        <button onClick={() => setCurrentView('websocket')}>Websocket</button>
      </div>

      {currentView === 'websocket' && <MyWebSocket />}
      {currentView === 'allUsers' && <Users />}
      {currentView === 'currentUser' && <User />}
    </>
  );
};

export default Dashboard;
