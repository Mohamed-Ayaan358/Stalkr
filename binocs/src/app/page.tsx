"use client";
import { useEffect, useState } from "react";
import { Website, connectWebSocket } from "./WebSocketExample";

export default function Home() {
  const [websocketData, setWebsocketData] = useState<Website[]>([]);

  useEffect(() => {
    const socket = connectWebSocket("ws://localhost:8080/ws", (data) => {
      // Update the state with the new WebSocket data
      setWebsocketData(data);
    });

    return () => {
      // Clean up the WebSocket connection when the component unmounts
      socket.close();
    };
  }, []);

  return (
    <main className="flex min-h-screen flex-col items-center justify-between p-24">
      <div className="z-10 w-full max-w-5xl items-center justify-between font-mono text-sm lg:flex">
        <p>WebSocket Data:</p>
        <ul>
          {websocketData.map((website, index) => (
            <li key={index}>
              <strong>ID:</strong> {website.id}
              <br></br>
              <strong>Name:</strong> {website.name}
              <br></br>
              <strong>URL:</strong> {website.url}
              <br></br>
              <strong>Hash:</strong> {website.hash}
              <br></br>
              <strong>Time:</strong> {website.time}
              <br></br>
              <strong>Changed:</strong> {website.changed.toString()}
              <br></br>
            </li>
          ))}
        </ul>
      </div>
    </main>
  );
}
