import React, { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";
import axios from "axios";

const App: React.FC = (props: React.PropsWithChildren<{}>) => {
  const [error, setError] = useState("");
  const history = useHistory();

  useEffect(() => {
    const id = setTimeout(() => {
      setError("");
    }, 3000);

    return () => clearTimeout(id);
  }, [error]);

  const handleSignOut = async () => {
    try {
      const response = await axios.get("api/signout");
      setError(response.data);
      history.push("/");
    } catch (err) {
      if (err.response.status === 401) {
        setError(err.response.data);

        setTimeout(() => {
          return history.push("/");
        }, 1500);
      }

      setError(err.response.data);
    }
  };

  return (
    <div className="app">
      <header className="app-header">
        <h1>Welcome to Go Starter!</h1>

        <div className="app-menu">
          <button onClick={handleSignOut}>Sign Out</button>
        </div>
      </header>

      {error ? <p className="error">{error}</p> : null}

      {props.children}
    </div>
  );
};

export default App;
