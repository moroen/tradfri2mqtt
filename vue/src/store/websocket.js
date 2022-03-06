var storeContext;

const wsURL = `ws://${window.location.host}/api/ws`;

const onConnect = (event) => {
  console.log("Successfully connected to the echo websocket server...");
};

const onMessage = (event) => {
  var obj = JSON.parse(event.data);
  storeContext.commit("addLogEntry", obj);
};

// Mutations
export const addLogEntry = (state, payload) => {
  state.websocket.log.push(payload);
};

export const clearLog = (state) => {
  state.websocket.log = [];
};

// Actions
export const sendWSCommand = (context, payload) => {
  const connection = context.state.websocket.connection;
  if (connection.readyState == WebSocket.OPEN) {
    connection.send(payload);
  } else {
    console.error("Unable to send WSCommand, socket not open");
  }
};

export const newWebSocket = (context) => {
  storeContext = context;
  const connection = new WebSocket(wsURL);
  connection.onopen = onConnect;
  connection.onmessage = onMessage;

  return connection;
};
