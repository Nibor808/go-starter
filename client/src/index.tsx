import "./index.css";
import React from "react";
import ReactDOM from "react-dom";
import App from "./App";
import Home from "./components/home";
import SignUpEmail from "./components/auth/sign_up_email";
import ConfirmEmail from "./components/auth/confirm_email";
import SignUpPassword from "./components/auth/sign_up_password";
import SignIn from "./components/auth/sign_in";
import Dashboard from "./components/dashboard";
import { BrowserRouter, Route, Switch } from "react-router-dom";
import reportWebVitals from "./reportWebVitals";

ReactDOM.render(
  <BrowserRouter>
    <React.StrictMode>
      <App>
        <Switch>
          <Route exact path="/" component={Home} />
          <Route exact path="/signupemail" component={SignUpEmail} />
          <Route path="/confirmemail/:token/:userID" component={ConfirmEmail} />
          <Route exact path="/signuppassword" component={SignUpPassword} />
          <Route exact path="/signin" component={SignIn} />
          <Route exact path="/dashboard" component={Dashboard} />
        </Switch>
      </App>
    </React.StrictMode>
  </BrowserRouter>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
