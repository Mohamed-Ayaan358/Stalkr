export interface Website {
  id: number;
  name: string;
  url: string;
  hash: string;
  time: number;
  changed: boolean;
}
let receivedWebsites: Website[] = [];

export const connectWebSocket = (
  url: string,
  onUpdateWebsites: (websites: Website[]) => void
) => {
  const socket = new WebSocket(url);

  socket.addEventListener("open", (event) => {
    console.log("WebSocket connection established:", event);
  });

  socket.addEventListener("message", (event) => {
    try {
      const data = JSON.parse(event.data);
      receivedWebsites = data;
      onUpdateWebsites(receivedWebsites);
    } catch (error) {
      console.error("Error parsing WebSocket message:", error);
    }
  });

  socket.addEventListener("close", (event) => {
    console.log("WebSocket connection closed:", event);
  });

  socket.addEventListener("error", (error) => {
    console.error("WebSocket error:", error);
  });

  return socket;
};
