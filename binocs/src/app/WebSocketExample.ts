export const connectWebSocket = (
  url: string,
  onMessage: (data: string) => void
) => {
  const socket = new WebSocket(url);

  socket.addEventListener("open", (event) => {
    console.log("WebSocket connection established:", event);
  });

  socket.addEventListener("message", (event) => {
    console.log("Received message:", event.data);
    onMessage(event.data);
  });

  socket.addEventListener("close", (event) => {
    console.log("WebSocket connection closed:", event);
  });

  socket.addEventListener("error", (error) => {
    console.error("WebSocket error:", error);
  });

  return socket;
};
