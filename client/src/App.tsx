import { useEffect, useState } from 'react';
import { useHistory } from 'react-router-dom';
import axios from 'axios';

const App: React.FC = props => {
  const [message, setMessage] = useState('');
  const history = useHistory();

  useEffect(() => {
    const id = setTimeout(() => {
      setMessage('');
    }, 3000);

    return () => clearTimeout(id);
  }, [message]);

  const handleSignOut = async () => {
    try {
      const response = await axios.get('api/signout');
      setMessage(response.data);
      history.push('/');
    } catch (err) {
      setMessage(err.response.data);
    }
  };

  return (
    <div className='app'>
      <header className='app-header'>
        <h1>Welcome to Go Starter!</h1>

        <div className='app-menu'>
          <button onClick={handleSignOut}>Sign Out</button>
        </div>
      </header>

      {message ? <p className='error'>{message}</p> : null}

      {props.children}
    </div>
  );
};

export default App;
