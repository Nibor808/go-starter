import React, { useState, useEffect, useRef } from "react";

const MyWebSocket: React.FC = () => {
  const socket = useRef<WebSocket | null>(null);
  const [message, setMessage] = useState("");
  const [error, setError] = useState("");
  const [websocketResponse, setWebsocketResponse] = useState("");
  const [counter, setCounter] = useState(0);

  useEffect(() => {
    const id = setTimeout(() => {
      setMessage("");
      setError("");
    }, 3000);

    return () => clearTimeout(id);
  }, [message]);

  useEffect(() => {
    socket.current = new WebSocket("ws://localhost:3000/wsconnect");
  }, []);

  useEffect(() => {
    if (socket.current) {
      socket.current.onmessage = (ev: MessageEvent) => {
        setWebsocketResponse(ev.data);
        setCounter(counter + 1);
      };
    }
  }, [counter]);

  useEffect(() => {
    if (socket.current) {
      socket.current.onopen = () => {
        setMessage("SOCKET OPENED");
      };

      socket.current.onclose = (ev: CloseEvent) => {
        setMessage("SOCKET CLOSED: " + ev.code);
      };

      socket.current.onerror = () => {
        setError("There was an error with the websocket.");
      };
    }

    return () => {
      socket.current?.close();
      setCounter(0);
    };
  }, []);

  useEffect(() => {}, []);

  const handleMessage = (ev: React.MouseEvent<HTMLButtonElement>) => {
    const elem = document.getElementById("messageInput") as HTMLInputElement;
    const msg: string = elem.value;
    ev.preventDefault();

    if (socket.current?.readyState === socket.current?.OPEN) {
      socket.current?.send(msg);
    } else {
      setError("WebSocket closed. Refresh to re-open");
      setWebsocketResponse("");
      setCounter(0);
      elem.value = "";
    }
  };

  return (
    <>
      <div className="websocket-header">
        <p>Websocket</p>

        {message ? <p className="success">{message}</p> : null}
      </div>

      {error ? <p className="error">{error}</p> : null}

      <p>Websocket Response: {websocketResponse}</p>
      <p>Responses Recieved: {counter}</p>

      <div className="v-form">
        <label htmlFor="messageInput">Message</label>
        <input type="text" id="messageInput" />
      </div>
      <button
        onClick={(ev: React.MouseEvent<HTMLButtonElement>) => handleMessage(ev)}
      >
        Send Message
      </button>
    </>
  );
};

export default MyWebSocket;
