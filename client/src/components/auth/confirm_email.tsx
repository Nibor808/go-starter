import React, { useState, useEffect } from "react";
import { useParams, useHistory } from "react-router-dom";
import axios from "axios";

interface RouteParams {
  token: string;
  userID: string;
}

const ConfirmEmail: React.FC = () => {
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const { token, userID }: RouteParams = useParams();
  const history = useHistory();

  useEffect(() => {
    (async () => {
      try {
        const response = await axios.get(
          `http://localhost:3000/api/confirmemaildata/${token}/${userID}`
        );
        setSuccess(response.data);
        history.push("/signuppassword");
      } catch (err) {
        setError(err.response.data);
      }
    })();
  }, []);

  return (
    <div>
      <h3>Confirm Email</h3>

      <p>{error}</p>
      <p>{success}</p>
    </div>
  );
};

export default ConfirmEmail;
