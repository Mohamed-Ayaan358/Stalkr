"use client";
import { useEffect, useState } from "react";
import { connectWebSocket } from "./WebSocketExample";

export default function Home() {
  const [websocketData, setWebsocketData] = useState<Set<string>>(new Set());

  useEffect(() => {
    const socket = connectWebSocket("ws://localhost:8080/ws", (data) => {
      // Update the state with the new WebSocket data
      setWebsocketData((prevData: any) => new Set([...prevData, data]));
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
          {Array.from(websocketData).map((data, index) => (
            <li key={index}>{data}</li>
          ))}
        </ul>
      </div>
    </main>
  );
}
