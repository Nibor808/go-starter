import React, { useState, useEffect } from "react";
import { useHistory } from "react-router-dom";
import axios from "axios";

const SignIn: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const history = useHistory();

  useEffect(() => {
    setTimeout(() => {
      setError("");
    }, 2000);
  }, [error]);

  const handleSubmit: React.FormEventHandler = async (ev: React.FormEvent) => {
    ev.preventDefault();

    const formData = new FormData();
    formData.append("email", email);
    formData.append("password", password);

    try {
      await axios
        .post("api/signin", formData)
        .then(() => history.push("/dashboard"));
    } catch (err) {
      setError(err.response.data);
    }
  };

  return (
    <div>
      <h3>Sign In</h3>

      {error ? <p className="error">{error}</p> : null}

      <form onSubmit={handleSubmit} className="v-form">
        <label htmlFor="email">email</label>
        <input
          id="email"
          type="email"
          onChange={(ev: React.ChangeEvent<HTMLInputElement>) =>
            setEmail(ev.target.value)
          }
        />

        <label htmlFor="password">password</label>
        <input
          id="password"
          type="password"
          onChange={(ev: React.ChangeEvent<HTMLInputElement>) =>
            setPassword(ev.target.value)
          }
        />

        <button type="submit">Sign in</button>
      </form>
    </div>
  );
};

export default SignIn;
