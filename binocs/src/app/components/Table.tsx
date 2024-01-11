"use client";
import { useEffect, useState } from "react";
import { Website, connectWebSocket } from "@/app/WebSocketExample";
import axios from "axios";

export default function TableContent() {
  const [websocketData, setWebsocketData] = useState<Website[]>([]);
  async function fetchData() {
    axios.get("/api/get").then((res) => {
      if (res.data.data.data) {
        console.log(res.data.data.data);
        setWebsocketData(res.data.data.data);
      }
    });
  }
  async function deleteData(name: string) {
    const response = await axios.post("/api/delete", { websiteName: name });
    location.reload();
  }

  useEffect(() => {
    const socket = connectWebSocket("ws://localhost:8080/ws", (data) => {
      setWebsocketData(data);
    });
    if (socket.readyState === WebSocket.OPEN) {
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
    <>
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
            {websocketData &&
              websocketData.map((website, index) => (
                <tr key={index} className="bg-gray-100">
                  <td
                    className="border px-2 py-2"
                    title={website.id.toString()}
                  >
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
                    <button
                      className="focus:outline-none text-white bg-red-700 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-red-600 dark:hover:bg-red-700 dark:focus:ring-red-900"
                      onClick={() => deleteData(website.name)}
                    >
                      {" "}
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
          </tbody>
        </table>
      </div>
    </>
  );
}
