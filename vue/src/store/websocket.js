var storeContext;
var connectionDelay = 5000;
const connectionMaxDelay = 25000;
const connectionDelayDelta = 5000;

import { showError } from "./actions";

const wsURL = `ws://${window.location.host}/api/ws`;

const onConnect = (event) => {
  console.log("Successfully connected to the echo websocket server...");
};

const onMessage = (event) => {
  // console.log(event.data);
  var obj = JSON.parse(event.data);
  // storeContext.commit("addWSLogEntry", obj);
  console.log(JSON.stringify(obj));
  switch (obj.class) {
    case "log":
      storeContext.commit("addLogEntry", obj.data);
      break;
    case "devices":
      storeContext.commit("updateDeviceInfo", obj.data);
      break;
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
  let connection = context.state.websocket.connection;

  if (connection == null) {
    storeContext = context;
    connection = new WebSocket(wsURL);
    connection.onmessage = onMessage;

    connection.onopen = () => {
      context.commit("setStatus", "WebSocket connected");
      connectionDelay = 5000;
      context.dispatch("sendWSCommand", payload);
    };

    connection.onclose = () => {
      context.commit("setStatus", "WebSocket not connected");
      console.log("Connection closed");
      context.commit("setConnection", null);
    };

    connection.onerror = () => {
      context.commit("setStatus", "WebSocket failed to connect");
      setTimeout(() => {
        context.dispatch("sendWSCommand", payload);
      }, connectionDelay);

      if (connectionDelay < connectionMaxDelay) {
        connectionDelay = connectionDelay + connectionDelayDelta;
      } else {
        connectionDelay = connectionMaxDelay;
      }
    };

    context.commit("setConnection", connection);

    return;
  }

  if (connection.readyState == WebSocket.OPEN) {
    connection.send(payload);
  } else {
    context.commit("setConnection", null);
    context.dispatch("sendWSCommand", payload);

    console.error("Unable to send WSCommand, socket not open");
  }
};
