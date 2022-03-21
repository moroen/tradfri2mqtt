var storeContext;

import { showError } from "./actions";

const wsURL = `ws://${window.location.host}/api/ws`;

const onConnect = (event) => {
  console.log("Successfully connected to the echo websocket server...");
};

const onMessage = (event) => {
  // console.log(event.data)
  var obj = JSON.parse(event.data);
  storeContext.commit("addWSLogEntry", obj);
  switch (obj.class) {
    case "log":
      storeContext.commit("addLogEntry", obj.data);
    case "devices":
      storeContext.commit("updateDeviceInfo", obj.data);
  }
};

// Mutations
export const addLogEntry = (state, payload) => {
  state.websocket.log.push(payload);
};

export const addWSLogEntry = (state, payload) => {
  state.websocket.wslog.push(payload);
};

export const updateDeviceInfo = (state, payload) => {
  state.websocket.devices[payload.id] = payload;
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
    showError("Unable to send WSCommand, socket not open");
  }
};

export const newWebSocket = (context) => {
  storeContext = context;
  const connection = new WebSocket(wsURL);
  connection.onopen = onConnect;
  connection.onmessage = onMessage;
  connection.onerror = () => {};

  return connection;
};
