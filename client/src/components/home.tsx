import React from 'react';
import { Link } from 'react-router-dom';

const Home: React.FC = () => {
    return (
        <>
            <h3>Sign in or sign up.</h3>

            <div className="app-sub-menu">
                <Link to="signupemail">create user</Link>

                <Link to="signin">sign in</Link>
            </div>
        </>
    );
};

export default Home;
