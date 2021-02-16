import { Link } from 'react-router-dom';

const Home: React.FC = () => {
  return (
    <div>
      <h3>Sign in or sign up.</h3>
      <Link to='signupemail'>sign up</Link>

      <Link to='signin'>sign in</Link>
    </div>
  );
};

export default Home;
