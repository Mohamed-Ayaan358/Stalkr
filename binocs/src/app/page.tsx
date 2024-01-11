"use client";
import { useEffect, useState } from "react";
import { Website, connectWebSocket } from "./WebSocketExample";
import ButtonModal from "./components/ButtonModal";
import axios from "axios";

export default function Home() {
  const [websocketData, setWebsocketData] = useState<Website[]>([]);
  async function fetchData() {
    axios.get("/api/get").then((res) => {
      if (res.data.data.data) {
        setWebsocketData(res.data.data.data);
      }
    });
  }
  function deleteData(name: string) {
    axios.post("/api/delete", { websiteName: name });
  }

  useEffect(() => {
    const socket = connectWebSocket("ws://localhost:8080/ws", (data) => {
      setWebsocketData(data);
    });
    if (socket.readyState === WebSocket.OPEN) {
      fetchData();
      return () => {
        socket.close();
      };
    } else {
      socket.addEventListener("open", () => {
        fetchData();
        return () => {
          socket.close();
        };
      });
    }
  }, []);

  function truncateText(text: any, maxLength: any) {
    if (text.length <= maxLength) {
      return text;
    } else {
      return text.substring(0, maxLength) + "...";
    }
  }

  // Keep the tooltip but maybe dont need the responsive-td class, You would need a truncate class maybe?
  return (
    <main className="flex flex-col items-center justify-between p-4 md:p-24">
      <p className="text-lg font-bold mb-4">WebSocket Data:</p>
      <div className="overflow-x-auto w-full">
        <table className="min-w-full table-auto">
          <thead>
            <tr>
              <th className=" py-2">ID</th>
              <th className=" py-2">Name</th>
              <th className=" py-2">URL</th>
              <th className=" py-2">Hash</th>
              <th className=" py-2">Time</th>
              <th className=" py-2">Changed</th>
            </tr>
          </thead>
          <tbody>
            {websocketData.map((website, index) => (
              <tr key={index} className="bg-gray-100">
                <td className="border px-2 py-2" title={website.id.toString()}>
                  {website.id}
                </td>
                <td className="border px-2 py-2" title={website.name}>
                  {website.name}
                </td>
                <td className="border px-2 py-2 " title={website.url}>
                  <a href={truncateText(website.url, 30)}>
                    {truncateText(website.url, 30)}
                  </a>
                </td>
                <td className="border px-2 py-2 " title={website.hash}>
                  {website.hash}
                </td>
                <td
                  className="border px-2 py-2"
                  title={website.time.toString()}
                >
                  {website.time}
                </td>
                <td
                  className="border px-2 py-2"
                  title={website.changed.toString()}
                >
                  {website.changed.toString()}
                </td>
                <td
                  className="border px-2 py-2"
                  title={website.changed.toString()}
                >
                  <button onClick={() => deleteData(website.name)}>
                    {" "}
                    Delete
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <ButtonModal />
    </main>
  );
}
