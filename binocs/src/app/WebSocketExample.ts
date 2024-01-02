export interface Website {
  id: number;
  name: string;
  url: string;
  hash: string;
  time: number;
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
    console.log("Received message:", event.data);

    // Assuming event.data is a JSON string representing an array of websites
    try {
      const data = JSON.parse(event.data);

      // Update the receivedWebsites variable
      receivedWebsites = data;

      // Call the provided callback to update the state or perform any action
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
