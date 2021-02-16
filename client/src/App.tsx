import { Link } from 'react-router-dom';

const App: React.FC = props => {
  return (
    <div className='App'>
      <header className='App-header'>
        <h1>Welcome to Go Starter!</h1>
        <Link to='/'>Home</Link>

        <div>{props.children}</div>
      </header>
    </div>
  );
};

export default App;
