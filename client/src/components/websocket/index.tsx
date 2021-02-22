import { useState, useEffect, useRef } from 'react';

const MyWebSocket: React.FC = () => {
  const socket = useRef<WebSocket | null>(null);
  const [message, setMessage] = useState('');
  const [websocketResponse, setWebsocketResponse] = useState('');
  const [counter, setCounter] = useState(0);

  useEffect(() => {
    setTimeout(() => {
      setMessage('');
    }, 2000);
  }, [message]);

  useEffect(() => {
    socket.current = new WebSocket('ws://localhost:3000/wsconnect');
  }, []);

  useEffect(() => {
    if (socket.current) {
      socket.current.onopen = () => {
        setMessage('SOCKET OPENED');
      };

      socket.current.onmessage = (ev: MessageEvent) => {
        setWebsocketResponse(ev.data);
        setCounter(counter + 1);
      };

      socket.current.onclose = (ev: CloseEvent) => {
        setMessage('SOCKET CLOSED: ' + ev.code);
      };

      socket.current.onerror = (ev: Event) => {
        setMessage('ERROR: ' + ev);
      };
    }
  }, [counter]);

  const handleMessage = (ev: React.MouseEvent<HTMLButtonElement>) => {
    const elem = document.getElementById('messageInput') as HTMLInputElement;
    const msg: string = elem.value;
    ev.preventDefault();

    if (socket.current?.readyState === socket.current?.OPEN) {
      socket.current?.send(msg);
    } else {
      setMessage('WebSocket closed. Refresh to re-open');
      setWebsocketResponse('');
      setCounter(0);
      elem.value = '';
    }
  };

  return (
    <>
      <p>Messages: {message}</p>

      <p>Websocket Response: {websocketResponse}</p>
      <p>Responses Recieved: {counter}</p>

      <div>
        <label htmlFor='messageInput'>Message</label>
        <input type='text' id='messageInput' />
      </div>
      <button onClick={ev => handleMessage(ev)}>Send Message</button>
    </>
  );
};

export default MyWebSocket;
