import React from 'react';
import ReactDOM from 'react-dom';
import App from './App';
import Home from './components/home';
import SignUp from './components/sign_up';
import { BrowserRouter, Route, Switch } from 'react-router-dom';
import reportWebVitals from './reportWebVitals';

ReactDOM.render(
  <BrowserRouter>
    <React.StrictMode>
      <App>
        <Switch>
          <Route exact path='/' component={Home} />
          <Route exact path='/signup' component={SignUp} />
        </Switch>
      </App>
    </React.StrictMode>
  </BrowserRouter>,
  document.getElementById('root')
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
