import { useState, useEffect } from 'react';
import { Link, useHistory } from 'react-router-dom';
import axios from 'axios';

const App: React.FC = props => {
  const [message, setMessage] = useState('');
  const history = useHistory();

  useEffect(() => {
    setTimeout(() => {
      setMessage('');
    }, 2000);
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
    <div className='App'>
      <header className='App-header'>
        <h1>Welcome to Go Starter!</h1>
        <p>{message}</p>
        <Link to='/'>Home</Link>
        <button onClick={handleSignOut}>Sign Out</button>

        <div>{props.children}</div>
      </header>
    </div>
  );
};

export default App;
