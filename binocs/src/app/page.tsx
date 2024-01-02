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
              <strong>ID:</strong> {website.id},<strong>Name:</strong>{" "}
              {website.name},<strong>URL:</strong> {website.url},
              <strong>Hash:</strong> {website.hash},<strong>Time:</strong>{" "}
              {website.time}
            </li>
          ))}
        </ul>
      </div>
    </main>
  );
}
