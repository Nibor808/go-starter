import React, { useState, useEffect } from "react";
import axios from "axios";
import { IUser } from "./interfaces";
import { useHistory } from "react-router-dom";

const User: React.FC = () => {
  const history = useHistory();
  const [error, setError] = useState("");
  const [user, setUser] = useState<IUser>({
    id: "",
    email: "",
    isActive: false,
    isAdmin: false,
  });

  useEffect(() => {
    const id = setTimeout(() => {
      setError("");
    }, 2000);

    return () => clearTimeout(id);
  }, [error]);

  useEffect(() => {
    (async () => {
      try {
        const response = await axios.get("api/user");

        setUser(response.data);
      } catch (err) {
        if (err.response.status === 401) {
          setError(err.response.data);

          setTimeout(() => {
            return history.push("/");
          }, 1500);
        }

        setError(err.response.data);
      }
    })();
  }, [history]);

  return (
    <>
      <p>User From Cookie</p>
      {error ? <p className="error">ERROR: {error}</p> : null}
      <ul>
        <li>id: {user.id}</li>
        <li>email: {user.email}</li>
        <li>isActive: {String(user.isActive)}</li>
        <li>isAdmin: {String(user.isAdmin)}</li>
      </ul>
    </>
  );
};

export default User;
