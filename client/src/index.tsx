import './index.css';
import React from 'react';
import ReactDOM from 'react-dom';
import { App } from './App';
import Home from './components/home';
import SignUpEmail from './components/auth/sign_up_email';
import ConfirmEmail from './components/auth/confirm_email';
import SignUpPassword from './components/auth/sign_up_password';
import SignIn from './components/auth/sign_in';
import Dashboard from './components/dashboard';
import { BrowserRouter, Route, Routes } from 'react-router-dom';
import reportWebVitals from './reportWebVitals';

ReactDOM.render(
    <BrowserRouter>
        <React.StrictMode>
            <App>
                <Routes>
                    <Route path="/" element={<Home />} />
                    <Route path="/signupemail" element={<SignUpEmail />} />
                    <Route path="/confirmemail/:token/:userID" element={<ConfirmEmail />} />
                    <Route path="/signuppassword" element={<SignUpPassword />} />
                    <Route path="/signin" element={<SignIn />} />
                    <Route path="/dashboard" element={<Dashboard />} />
                </Routes>
            </App>
        </React.StrictMode>
    </BrowserRouter>,
    document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
